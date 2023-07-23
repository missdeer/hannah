#include "netease.h"

namespace ProviderAPI {
    namespace netease {

        boost::beast::http::response<boost::beast::http::string_body> playlist(std::string &reqBody, UrlQuery &urlQuery, UrlParam &params)
        {
            return {};
        }
        boost::beast::http::response<boost::beast::http::string_body> songinfo(std::string &reqBody, UrlQuery &urlQuery, UrlParam &params)
        {
            return {};
        }
    } // namespace netease
} // namespace ProviderAPI