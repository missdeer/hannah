#pragma once

#include <string>
#include <boost/beast/http.hpp>

class UrlQuery;
class UrlParam;

namespace ProviderAPI
{
    namespace qq
    {
        boost::beast::http::response<boost::beast::http::string_body> playlist(std::string &reqBody, UrlQuery &urlQuery, UrlParam &params);
        boost::beast::http::response<boost::beast::http::string_body> songinfo(std::string &reqBody, UrlQuery &urlQuery, UrlParam &params);
    } // namespace qq
} // namespace ProviderAPI