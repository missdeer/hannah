#pragma once

#include <memory>

#include <boost/asio/ip/tcp.hpp>
#include <boost/beast/core.hpp>

class BeastServer : public std::enable_shared_from_this<BeastServer>
{
    boost::asio::io_context       &m_ioc;
    boost::asio::ip::tcp::acceptor m_acceptor;

public:
    BeastServer(boost::asio::io_context &ioc, boost::asio::ip::tcp::endpoint endpoint);

    // Start accepting incoming connections
    void run();

private:
    void do_accept();

    void on_accept(boost::beast::error_code ec, boost::asio::ip::tcp::socket socket);

    void fail(boost::beast::error_code ec, char const *what);
};

void StartBeastServer();

void StopBeastServer();

void SetListenPort(unsigned short port);