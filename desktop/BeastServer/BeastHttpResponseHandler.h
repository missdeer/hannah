#pragma once

#include <boost/beast/core.hpp>
#include <boost/beast/http.hpp>

namespace beast = boost::beast; // from <boost/beast.hpp>
namespace http  = beast::http;  // from <boost/beast/http.hpp>

namespace BeastHttpServer
{
    // Returns a bad request response
    template<class Body, class Allocator>
    http::response<http::string_body> bad_request(http::request<Body, http::basic_fields<Allocator>> &req, beast::string_view why)
    {
        http::response<http::string_body> res {http::status::bad_request, req.version()};
        res.set(http::field::server, BOOST_BEAST_VERSION_STRING);
        res.set(http::field::content_type, "text/html; charset=UTF-8");
        res.keep_alive(req.keep_alive());
        res.body() = std::string(why);
        res.prepare_payload();
        return res;
    };

    // Returns a not found response
    template<class Body, class Allocator>
    http::response<http::string_body> not_found(http::request<Body, http::basic_fields<Allocator>> &req, beast::string_view target)
    {
        http::response<http::string_body> res {http::status::not_found, req.version()};
        res.set(http::field::server, BOOST_BEAST_VERSION_STRING);
        res.set(http::field::content_type, "text/html; charset=UTF-8");
        res.keep_alive(req.keep_alive());
        res.body() = "The resource '" + std::string(target) + "' was not found.";
        res.prepare_payload();
        return res;
    };

    // Returns a server error response
    template<class Body, class Allocator>
    http::response<http::string_body> server_error(http::request<Body, http::basic_fields<Allocator>> &req, beast::string_view what)
    {
        http::response<http::string_body> res {http::status::internal_server_error, req.version()};
        res.set(http::field::server, BOOST_BEAST_VERSION_STRING);
        res.set(http::field::content_type, "text/html; charset=UTF-8");
        res.keep_alive(req.keep_alive());
        res.body() = "An error occurred: '" + std::string(what) + "'";
        res.prepare_payload();
        return res;
    };

    http::response<http::string_body> ok(const std::string &resp_body, const std::string &content_type = "application/json; charset=UTF-8");
    http::response<http::string_body> htmlOk(const std::string &resp_body);
    http::response<http::string_body> jsonOk(const std::string &resp_body);
} // namespace BeastHttpServer