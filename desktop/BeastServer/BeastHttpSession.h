#pragma once

#include <memory>
#include <vector>

#include <boost/beast/core.hpp>
#include <boost/beast/http.hpp>

// Handles an HTTP server connection
class BeastHttpSession : public std::enable_shared_from_this<BeastHttpSession>
{
    // This queue is used for HTTP pipelining.
    class queue
    {
        enum
        {
            // Maximum number of responses we will queue
            limit = 8
        };

        // The type-erased, saved work item
        struct work
        {
            virtual ~work()           = default;
            virtual void operator()() = 0;
        };

        BeastHttpSession                  &self_;
        std::vector<std::unique_ptr<work>> items_;

    public:
        explicit queue(BeastHttpSession &self);

        // Returns `true` if we have reached the queue limit
        bool is_full() const;

        // Called when a message finishes sending
        // Returns `true` if the caller should initiate a read
        bool on_write();

        // Called by the HTTP handler to send a response.
        template<bool isRequest, class Body, class Fields>
        void operator()(boost::beast::http::message<isRequest, Body, Fields> &&msg)
        {
            // This holds a work item
            struct work_impl : work
            {
                BeastHttpSession                                    &self_;
                boost::beast::http::message<isRequest, Body, Fields> msg_;

                work_impl(BeastHttpSession &self, boost::beast::http::message<isRequest, Body, Fields> &&msg)
                    : self_(self), msg_(std::move(msg))
                {
                }

                void operator()()
                {
                    boost::beast::http::async_write(self_.stream_,
                                                    msg_,
                                                    boost::beast::bind_front_handler(&BeastHttpSession::on_write,
                                                                                     self_.shared_from_this(),
                                                                                     msg_.need_eof()));
                }
            };

            // Allocate and store the work
            items_.push_back(boost::make_unique<work_impl>(self_, std::move(msg)));

            // If there was no previous work, start this one
            if (items_.size() == 1)
                (*items_.front())();
        }
    };

    boost::beast::tcp_stream  stream_;
    boost::beast::flat_buffer buffer_;
    queue                     queue_;

    // The parser is stored in an optional container so we can
    // construct it from scratch it at the beginning of each new message.
    boost::optional<boost::beast::http::request_parser<boost::beast::http::string_body>> parser_;

public:
    // Take ownership of the socket
    explicit BeastHttpSession(boost::asio::ip::tcp::socket &&socket);

    // Start the session
    void run();

private:
    void do_read();

    void on_read(boost::beast::error_code ec, std::size_t bytes_transferred);

    void on_write(bool close, boost::beast::error_code ec, std::size_t bytes_transferred);

    void do_close();

    void fail(boost::beast::error_code ec, char const *what);
};
