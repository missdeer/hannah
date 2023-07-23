#include "kugou.h"

namespace ProviderAPI {
    namespace kugou {

        boost::beast::http::response<boost::beast::http::string_body> playlist(std::string &reqBody, UrlQuery &urlQuery, UrlParam &params)
        {
            return {};
        }
        boost::beast::http::response<boost::beast::http::string_body> songinfo(std::string &reqBody, UrlQuery &urlQuery, UrlParam &params)
        {
            return {};
        }
    } // namespace kugou
} // namespace ProviderAPI