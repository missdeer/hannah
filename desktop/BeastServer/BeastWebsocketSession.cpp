#include <algorithm>
#include <cstdlib>
#include <functional>
#include <iostream>
#include <string>
#include <thread>
#include <vector>
#include <boost/algorithm/string.hpp>
#include <boost/asio/bind_executor.hpp>
#include <boost/asio/dispatch.hpp>
#include <boost/asio/signal_set.hpp>
#include <boost/asio/strand.hpp>
#include <boost/beast/http.hpp>
#include <boost/beast/version.hpp>
#include <boost/make_unique.hpp>
#include <boost/optional.hpp>

#include "BeastWebsocketSession.h"
#include "BeastWebsocketSessionManager.h"

namespace beast     = boost::beast;         // from <boost/beast.hpp>
namespace http      = beast::http;          // from <boost/beast/http.hpp>
namespace websocket = beast::websocket;     // from <boost/beast/websocket.hpp>
namespace net       = boost::asio;          // from <boost/asio.hpp>
using tcp           = boost::asio::ip::tcp; // from <boost/asio/ip/tcp.hpp>

BeastWebsocketSession::BeastWebsocketSession(tcp::socket &&socket) : m_websocket(std::move(socket)) {}

BeastWebsocketSession::~BeastWebsocketSession()
{
    BeastWebsocketSessionManager::instance().remove(this);
}

void BeastWebsocketSession::onAccept(beast::error_code ec)
{
    if (ec)
    {
        fail(ec, "accept");
        return;
    }

    BeastWebsocketSessionManager::instance().add(this);
    // always an empty session id
    // set preset command by session id

    // Read a message
    m_websocket.async_read(m_buffer, beast::bind_front_handler(&BeastWebsocketSession::onRead, shared_from_this()));
}

void BeastWebsocketSession::onRead(beast::error_code ec, std::size_t bytes_transferred)
{
    boost::ignore_unused(bytes_transferred);

    if (ec)
    {
        fail(ec, "read");
        return;
    }

    // pass to handlers
    if (m_websocket.got_text())
    {
        std::string payload(beast::buffers_to_string(m_buffer.data()));

        std::string str1 = "on_message called with hdl: ?? and message : " + payload;

        try
        {
            bool bHandled = analyseCommand(payload);

            if (!bHandled)
            {
                onWebMsg(payload);

                //! NOTICE: it's mainly a test purpose work flow
                // Send to all connections
                // BeastWebsocketSessionManager::instance().send(payload);
            }
        }
        catch (const std::exception &e)
        {
            BeastWebsocketSessionManager::instance().send(e.what());
        }
        catch (...)
        {
            BeastWebsocketSessionManager::instance().send("unknown error on analyzing command and handle web message");
        }
    }

    // Clear the buffer
    m_buffer.consume(m_buffer.size());
    // wait for the next message
    m_websocket.async_read(m_buffer, beast::bind_front_handler(&BeastWebsocketSession::onRead, shared_from_this()));
}

void BeastWebsocketSession::onWrite(beast::error_code ec, std::size_t bytes_transferred)
{
    boost::ignore_unused(bytes_transferred);

    // Handle the error, if any
    if (ec)
    {
        fail(ec, "write");
        return;
    }

    // Remove the string from the queue
    m_msgQueue.erase(m_msgQueue.begin());

    // Send the next message if any
    if (!m_msgQueue.empty())
    {
        m_websocket.async_write(net::buffer(*m_msgQueue.front()), beast::bind_front_handler(&BeastWebsocketSession::onWrite, shared_from_this()));
    }
}

void BeastWebsocketSession::onSend(const std::shared_ptr<const std::string> &msg)
{ // Always add to queue
    m_msgQueue.push_back(msg);

    // Are we already writing?
    if (m_msgQueue.size() > 1)
    {
        return;
    }

    // We are not currently writing, so send this immediately
    m_websocket.async_write(net::buffer(*m_msgQueue.front()), beast::bind_front_handler(&BeastWebsocketSession::onWrite, shared_from_this()));
}

void BeastWebsocketSession::onWebMsg(const std::string &msg)
{
    BeastWebsocketSessionManager::instance().onWebMsg(msg);
}

void BeastWebsocketSession::fail(beast::error_code ec, char const *what)
{
    // Don't report these
    if (ec == net::error::operation_aborted || ec == websocket::error::closed)
    {
        return;
    }
}

std::string BeastWebsocketSession::getPresetCommand()
{
    std::string strCommand = getCommand("SESSION_1_DOWN");

    return strCommand;
}

std::string BeastWebsocketSession::getConfirmCommand()
{
    std::string strCommand = getCommand("SESSION_2_DOWN");

    return strCommand;
}

std::string BeastWebsocketSession::getCommand(const std::string &actionType)
{
    return {};
}

bool BeastWebsocketSession::analyseCommand(const std::string &strCommand)
{
    return false;
}

void BeastWebsocketSession::send(const std::shared_ptr<const std::string> &msg)
{
    // Post our work to the strand, this ensures
    // that the members of `this` will not be
    // accessed concurrently.

    net::post(m_websocket.get_executor(), beast::bind_front_handler(&BeastWebsocketSession::onSend, shared_from_this(), msg));
}