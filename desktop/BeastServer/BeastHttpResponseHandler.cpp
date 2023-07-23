

#include <boost/beast/version.hpp>

#include "BeastHttpResponseHandler.h"

namespace BeastHttpServer
{
    http::response<http::string_body> ok(const std::string &resp_body, const std::string &content_type /*= "application/json; charset=UTF-8"*/)
    {
        auto const size = resp_body.size();

        http::response<http::string_body> res {http::status::ok, 11};
        res.set(http::field::server, BOOST_BEAST_VERSION_STRING);
        res.set(http::field::content_type, content_type);
        res.keep_alive(false);
        res.body() = resp_body;
        res.content_length(size);
        res.prepare_payload();
        return res;
    }

    http::response<http::string_body> htmlOk(const std::string &resp_body)
    {
        return ok(resp_body, "text/html; charset=UTF-8");
    }

    http::response<http::string_body> jsonOk(const std::string &resp_body)
    {
        return ok(resp_body, "application/json; charset=UTF-8");
    }

} // namespace BeastHttpServer