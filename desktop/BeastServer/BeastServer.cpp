

#include <algorithm>
#include <cstdlib>
#include <functional>
#include <iostream>
#include <string>
#include <thread>
#include <vector>
#include <boost/asio/bind_executor.hpp>
#include <boost/asio/dispatch.hpp>
#include <boost/asio/signal_set.hpp>
#include <boost/asio/strand.hpp>
#include <boost/beast/version.hpp>
#include <boost/format.hpp>
#include <boost/json.hpp>
#include <boost/make_unique.hpp>
#include <boost/optional.hpp>

#include "BeastServer.h"
#include "BeastHttpSession.h"
#include "BeastWebsocketSessionManager.h"

namespace beast = boost::beast;         // from <boost/beast.hpp>
namespace net   = boost::asio;          // from <boost/asio.hpp>
using tcp       = boost::asio::ip::tcp; // from <boost/asio/ip/tcp.hpp>

namespace
{
    net::io_context *g_ioc           = nullptr;
    unsigned short   g_listeningPort = 9100;
} // namespace

std::string ConcatAPIUrl(const std::string &path)
{
    if (!g_ioc)
    {
        return {};
    }
    boost::format fmter("http://127.0.0.1:%1%");
    fmter % g_listeningPort;
    return fmter.str() + path;
}

void StartBeastServer()
{
    const auto      address     = net::ip::make_address("127.0.0.1");
    const auto      threadCount = std::max<int>(4, std::thread::hardware_concurrency());
    net::io_context ioc {threadCount};
    g_ioc = &ioc;
    std::make_shared<BeastServer>(ioc, tcp::endpoint {address, g_listeningPort})->run();

    std::vector<std::thread> threads;
    threads.reserve(threadCount - 1);
    for (auto i = threadCount - 1; i > 0; --i)
    {
        threads.emplace_back([&ioc] {
            try
            {
                ioc.run();
            }
            catch (std::exception &e)
            {
            }
        });
    }

    boost::format fmter(R"({"pid" : %1%,"action" : "notification","type" : "websocket-ready","detail" : %2%})");
    DWORD         processId = GetCurrentProcessId();
    fmter % processId % g_listeningPort;
    std::string notification = fmter.str();

    try
    {
        ioc.run();
    }
    catch (std::exception &e)
    {
    }

    // (If we get here, it means StopBeastServer() is called)

    // Block until all the threads exit
    for (auto &thread : threads)
    {
        thread.join();
    }

    g_ioc = nullptr;

    BeastWebsocketSessionManager::instance().clear();
}

void StopBeastServer()
{
    if (g_ioc)
    {
        g_ioc->stop();
    }
}

BeastServer::BeastServer(boost::asio::io_context &ioc, boost::asio::ip::tcp::endpoint endpoint)
    : m_ioc(ioc), m_acceptor(boost::asio::make_strand(ioc))
{
    beast::error_code ec;

    // Open the acceptor
    m_acceptor.open(endpoint.protocol(), ec);
    if (ec)
    {
        fail(ec, "open");
        return;
    }

    // Allow address reuse
    m_acceptor.set_option(net::socket_base::reuse_address(true), ec);
    if (ec)
    {
        fail(ec, "set_option");
        return;
    }

    // Bind to the server address
    m_acceptor.bind(endpoint, ec);
    if (ec)
    {
        fail(ec, "bind");
        return;
    }

    // Start listening for connections
    m_acceptor.listen(net::socket_base::max_listen_connections, ec);
    if (ec)
    {
        fail(ec, "listen");
        return;
    }
}

void BeastServer::run()
{
    // We need to be executing within a strand to perform async operations
    // on the I/O objects in this session. Although not strictly necessary
    // for single-threaded contexts, this example code is written to be
    // thread-safe by default.
    boost::asio::dispatch(m_acceptor.get_executor(), beast::bind_front_handler(&BeastServer::do_accept, shared_from_this()));
}

void BeastServer::do_accept()
{
    // The new connection gets its own strand
    m_acceptor.async_accept(boost::asio::make_strand(m_ioc), beast::bind_front_handler(&BeastServer::on_accept, shared_from_this()));
}

void BeastServer::on_accept(beast::error_code ec, tcp::socket socket)
{
    if (ec)
    {
        fail(ec, "accept");
    }
    else
    {
        // Create the http session and run it
        std::make_shared<BeastHttpSession>(std::move(socket))->run();
    }

    // Accept another connection
    do_accept();
}

void BeastServer::fail(beast::error_code ec, char const *what)
{
    // Don't report on canceled operations
    if (ec == net::error::operation_aborted)
    {
        return;
    }
}

void SetListenPort(unsigned short port)
{
    g_listeningPort = port;
}