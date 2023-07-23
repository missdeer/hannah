

#include <functional>

#include "BeastWebsocketSessionManager.h"

BeastWebsocketSessionManager BeastWebsocketSessionManager::m_instance;

BeastWebsocketSessionManager::BeastWebsocketSessionManager()
{
    m_webMsgHandler = std::bind(&BeastWebsocketSessionManager::dummyWebMsgHandler, this, std::placeholders::_1);
}

void BeastWebsocketSessionManager::dummyWebMsgHandler(const std::string &msg) {}

void BeastWebsocketSessionManager::add(BeastWebsocketSessionPtr session)
{
    std::lock_guard<std::mutex> lock(m_beastWebsocketSessionsMutex);
    m_beastWebsocketSessions.insert(session);
}

void BeastWebsocketSessionManager::remove(BeastWebsocketSessionPtr session)
{
    std::lock_guard<std::mutex> lock(m_beastWebsocketSessionsMutex);
    m_beastWebsocketSessions.erase(session);
}

void BeastWebsocketSessionManager::clear()
{
    std::lock_guard<std::mutex> lock(m_beastWebsocketSessionsMutex);
    m_beastWebsocketSessions.clear();
}

BeastWebsocketSessions::iterator BeastWebsocketSessionManager::begin()
{
    return m_beastWebsocketSessions.begin();
}

BeastWebsocketSessions::iterator BeastWebsocketSessionManager::end()
{
    return m_beastWebsocketSessions.end();
}

void BeastWebsocketSessionManager::lock()
{
    m_beastWebsocketSessionsMutex.lock();
}

void BeastWebsocketSessionManager::unlock()
{
    m_beastWebsocketSessionsMutex.unlock();
}

void BeastWebsocketSessionManager::registerWebMsgHandler(std::function<void(const std::string &msg)> handler)
{
    m_webMsgHandler = handler;
}

void BeastWebsocketSessionManager::onWebMsg(const std::string &msg)
{
    m_webMsgHandler(msg);
}

void BeastWebsocketSessionManager::send(std::string message)
{
    // Put the message in a shared pointer so we can re-use it for each client
    const auto ss = std::make_shared<const std::string>(std::move(message));

    // Make a local list of all the weak pointers representing
    // the sessions, so we can do the actual sending without
    // holding the mutex:
    std::vector<std::weak_ptr<BeastWebsocketSession>> v;
    {
        std::lock_guard<std::mutex> lock(m_beastWebsocketSessionsMutex);
        v.reserve(m_beastWebsocketSessions.size());
        for (auto p : m_beastWebsocketSessions)
        {
            v.emplace_back(p->weak_from_this());
        }
    }

    // For each session in our local list, try to acquire a strong
    // pointer. If successful, then send the message on that session.
    for (auto const &wp : v)
    {
        if (auto sp = wp.lock())
        {
            sp->send(ss);
        }
    }
}
