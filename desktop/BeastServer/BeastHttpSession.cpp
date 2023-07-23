

#include <algorithm>
#include <cstdlib>
#include <functional>
#include <iostream>
#include <map>
#include <regex>
#include <string>
#include <vector>
#include <boost/asio/bind_executor.hpp>
#include <boost/asio/dispatch.hpp>
#include <boost/asio/signal_set.hpp>
#include <boost/asio/strand.hpp>
#include <boost/beast/version.hpp>
#include <boost/beast/websocket.hpp>
#include <boost/make_unique.hpp>
#include <boost/multi_index/identity.hpp>
#include <boost/multi_index/member.hpp>
#include <boost/multi_index/ordered_index.hpp>
#include <boost/multi_index_container.hpp>

#include "BeastHttpSession.h"
#include "BeastHttpResponseHandler.h"
#include "BeastWebsocketSessionManager.h"
#include "UrlQuery.h"

namespace beast     = boost::beast;         // from <boost/beast.hpp>
namespace http      = beast::http;          // from <boost/beast/http.hpp>
namespace websocket = beast::websocket;     // from <boost/beast/websocket.hpp>
namespace net       = boost::asio;          // from <boost/asio.hpp>
using tcp           = boost::asio::ip::tcp; // from <boost/asio/ip/tcp.hpp>

namespace
{
    using pHttpStaticRouterHandler = http::response<http::string_body> (*)(std::string &, UrlQuery &);
    struct HandlerDef
    {
        beast::string_view       path;
        http::verb               method;
        pHttpStaticRouterHandler handler;
        HandlerDef(const beast::string_view &p, http::verb m, pHttpStaticRouterHandler h) : path(p), method(m), handler(h) {}

        bool operator<(const HandlerDef &e) const
        {
            return path < e.path;
        }
    };

    using HandlerDefSet = boost::multi_index_container<
        HandlerDef,
        boost::multi_index::indexed_by<
            boost::multi_index::ordered_non_unique<boost::multi_index::member<HandlerDef, beast::string_view, &HandlerDef::path>>,
            boost::multi_index::ordered_non_unique<boost::multi_index::member<HandlerDef, http::verb, &HandlerDef::method>>>>;

    // http://172.16.206.13:8090/pages/viewpage.action?pageId=16867802
    static HandlerDefSet handlerMap = {
        //  {"/api/v1/database/create", http::verb::post, RESTfulAPI::database::post_create},
    };

    using pHttpRegexRouterHandler = http::response<http::string_body> (*)(std::string &, UrlQuery &, UrlParam &);
    struct RegexHandlerDef
    {
        std::regex              pattern;
        http::verb              method;
        pHttpRegexRouterHandler handler;
        RegexHandlerDef(const std::regex &p, http::verb m, pHttpRegexRouterHandler h) : pattern(p), method(m), handler(h) {}
    };

    using RegexHandlerDefSet                  = std::list<RegexHandlerDef>;
    static RegexHandlerDefSet regexHandlerMap = {
        //   {std::regex("^/api/v1/newView/tree/([a-zA-Z0-9_]+)$"), http::verb::put, RESTfulAPI::newview_tree::put},
    };
} // namespace

void registerHttpStaticRouter(const beast::string_view &path, http::verb method, pHttpStaticRouterHandler handler)
{
    handlerMap.insert({path, method, handler});
}

void registerHttpRegexRouter(const std::string &pattern, http::verb method, pHttpRegexRouterHandler handler)
{
    regexHandlerMap.emplace_back(std::regex(pattern), method, handler);
}

void registerHttpRegexRouter(const std::regex &regexp, http::verb method, pHttpRegexRouterHandler handler)
{
    regexHandlerMap.emplace_back(regexp, method, handler);
}

// This function produces an HTTP response for the given
// request. The type of the response object depends on the
// contents of the request, so the interface requires the
// caller to pass a generic lambda for receiving the response.
template<class Body, class Allocator, class Send> void handle_request(http::request<Body, http::basic_fields<Allocator>> &&req, Send &&send)
{
    // Make sure we can handle the method
    if (req.method() != http::verb::get && req.method() != http::verb::head && req.method() != http::verb::delete_ &&
        req.method() != http::verb::post && req.method() != http::verb::put)
    {
        return send(Cube::BeastHttpServer::bad_request(req, "Unknown HTTP-method"));
    }

    auto target = req.target();
    // Request path must be absolute and not contain "..".
    if (target.empty() || target[0] != '/')
    {
        return send(Cube::BeastHttpServer::bad_request(req, "Illegal request-target"));
    }
    std::string path(target);
    if (auto pos = path.find('?'); pos != std::string::npos)
    {
        path = path.substr(0, pos);
    }
    boost::beast::string_view trimmedPath(path);

    auto method        = req.method();
    auto [begin, end]  = handlerMap.get<0>().equal_range(trimmedPath);
    auto plaintextIter = std::find_if(begin, end, [method](const auto &item) { return item.method == method; });
    if (end != plaintextIter)
    {
        UrlQuery uriParam {std::string {target}};
        auto     res = plaintextIter->handler(req.body(), uriParam);
        return send(std::move(res));
    }

    auto regexIter = std::find_if(regexHandlerMap.begin(), regexHandlerMap.end(), [path](const auto &item) {
        return std::regex_match(path.begin(), path.end(), item.pattern);
    });
    if (regexHandlerMap.end() != regexIter)
    {
        UrlQuery           urlQuery {std::string {target}};
        std::smatch        paramMatch;
        UrlParam           urlParams;
        const std::string &urlPath(path);
        if (std::regex_search(urlPath, paramMatch, regexIter->pattern))
        {
            for (size_t i = 1, count = paramMatch.size(); i < count; ++i)
            {
                urlParams.push_back(paramMatch[i]);
            }
        }
        auto res = regexIter->handler(req.body(), urlQuery, urlParams);
        return send(std::move(res));
    }

    return send(Cube::BeastHttpServer::not_found(req, target));
}

BeastHttpSession::queue::queue(BeastHttpSession &self) : self_(self)
{
    static_assert(limit > 0, "queue limit must be positive");
    items_.reserve(limit);
}

bool BeastHttpSession::queue::is_full() const
{
    return items_.size() >= limit;
}

bool BeastHttpSession::queue::on_write()
{
    BOOST_ASSERT(!items_.empty());
    auto const was_full = is_full();
    items_.erase(items_.begin());
    if (!items_.empty())
    {
        (*items_.front())();
    }
    return was_full;
}

BeastHttpSession::BeastHttpSession(tcp::socket &&socket) : stream_(std::move(socket)), queue_(*this) {}

void BeastHttpSession::run()
{
    // We need to be executing within a strand to perform async operations
    // on the I/O objects in this session. Although not strictly necessary
    // for single-threaded contexts, this example code is written to be
    // thread-safe by default.
    net::dispatch(stream_.get_executor(), beast::bind_front_handler(&BeastHttpSession::do_read, this->shared_from_this()));
}

void BeastHttpSession::do_read()
{
    // Construct a new parser for each message
    parser_.emplace();

    // Apply a reasonable limit to the allowed size
    // of the body in bytes to prevent abuse.
    parser_->body_limit(50 * 1024 * 1024);

    // Set the timeout.
    stream_.expires_after(std::chrono::seconds(30));

    // Read a request using the parser-oriented interface
    http::async_read(stream_, buffer_, *parser_, beast::bind_front_handler(&BeastHttpSession::on_read, shared_from_this()));
}

void BeastHttpSession::on_read(beast::error_code ec, std::size_t bytes_transferred)
{
    boost::ignore_unused(bytes_transferred);

    // This means they closed the connection
    if (ec == http::error::end_of_stream)
    {
        return do_close();
    }

    if (ec)
    {
        fail(ec, "read");
        return;
    }

    // See if it is a WebSocket Upgrade
    if (websocket::is_upgrade(parser_->get()))
    {
        // Create a websocket session, transferring ownership
        // of both the socket and the HTTP request.
        std::make_shared<BeastWebsocketSession>(stream_.release_socket())->doAccept(parser_->release());
        return;
    }

    // Send the response
    handle_request(parser_->release(), queue_);

    // If we aren't at the queue limit, try to pipeline another request
    if (!queue_.is_full())
    {
        do_read();
    }
}

void BeastHttpSession::on_write(bool close, beast::error_code ec, std::size_t bytes_transferred)
{
    boost::ignore_unused(bytes_transferred);

    if (ec)
    {
        fail(ec, "write");
        return;
    }

    if (close)
    {
        // This means we should close the connection, usually because
        // the response indicated the "Connection: close" semantic.
        return do_close();
    }

    // Inform the queue that a write completed
    if (queue_.on_write())
    {
        // Read another request
        do_read();
    }
}

void BeastHttpSession::do_close()
{
    // Send a TCP shutdown
    beast::error_code ec;
    stream_.socket().shutdown(tcp::socket::shutdown_send, ec);

    // At this point the connection is closed gracefully
}

void BeastHttpSession::fail(beast::error_code ec, char const *what)
{
    // Don't report on canceled operations
    if (ec == net::error::operation_aborted)
    {
        return;
    }
}
