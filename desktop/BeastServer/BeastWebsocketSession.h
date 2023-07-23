#pragma once

#include <memory>
#include <mutex>
#include <string>
#include <unordered_set>
#include <vector>

#include <boost/asio/ip/tcp.hpp>
#include <boost/beast/core.hpp>
#include <boost/beast/http.hpp>
#include <boost/beast/websocket.hpp>

// Echoes back all received WebSocket messages
class BeastWebsocketSession : public std::enable_shared_from_this<BeastWebsocketSession>
{
    boost::beast::websocket::stream<boost::beast::tcp_stream> m_websocket;
    boost::beast::flat_buffer                                 m_buffer;
    std::vector<std::shared_ptr<const std::string>>           m_msgQueue;
    std::mutex                                                m_msgQueueMutex;
    std::string                                               m_strSessionID;
    std::string                                               m_strSessionToken;

public:
    // Take ownership of the socket
    explicit BeastWebsocketSession(boost::asio::ip::tcp::socket &&socket);

    ~BeastWebsocketSession();

    // Start the asynchronous accept operation
    template<class Body, class Allocator>
    void doAccept(boost::beast::http::request<Body, boost::beast::http::basic_fields<Allocator>> req)
    {
        // Set suggested timeout settings for the websocket
        m_websocket.set_option(
            boost::beast::websocket::stream_base::timeout::suggested(boost::beast::role_type::server));

        // Set a decorator to change the Server of the handshake
        m_websocket.set_option(
            boost::beast::websocket::stream_base::decorator([](boost::beast::websocket::response_type &res) {
                res.set(boost::beast::http::field::server,
                        std::string(BOOST_BEAST_VERSION_STRING) + " hannah-api-fake-server");
            }));
        m_websocket.auto_fragment(true);
        m_websocket.write_buffer_bytes(20 * 1024 * 1024);
        // Accept the websocket handshake
        m_websocket.async_accept(
            req, boost::beast::bind_front_handler(&BeastWebsocketSession::onAccept, shared_from_this()));
    }

    void send(const std::shared_ptr<const std::string> &msg);

private:
    void onAccept(boost::beast::error_code ec);

    void onRead(boost::beast::error_code ec, std::size_t bytes_transferred);

    void onWrite(boost::beast::error_code ec, std::size_t bytes_transferred);

    void onSend(const std::shared_ptr<const std::string> &msg);

    void onWebMsg(const std::string &msg);

    void fail(boost::beast::error_code ec, char const *what);

    std::string getPresetCommand();

    std::string getConfirmCommand();

    bool analyseCommand(const std::string &strCommand);

    std::string getCommand(const std::string &actionType);
};
