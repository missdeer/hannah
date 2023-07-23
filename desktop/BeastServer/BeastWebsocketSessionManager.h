#pragma once

#include "BeastWebsocketSession.h"

using BeastWebsocketSessionPtr = BeastWebsocketSession *;
using BeastWebsocketSessions   = std::unordered_set<BeastWebsocketSessionPtr>;

class BeastWebsocketSessionManager
{
public:
    static BeastWebsocketSessionManager &instance()
    {
        return m_instance;
    }

    void add(BeastWebsocketSessionPtr session);
    void remove(BeastWebsocketSessionPtr session);
    void clear();

    BeastWebsocketSessions::iterator begin();
    BeastWebsocketSessions::iterator end();
    void                             lock();
    void                             unlock();

    void registerWebMsgHandler(std::function<void(const std::string &msg)> handler);

    void onWebMsg(const std::string &msg);

    void send(std::string message);

private:
    static BeastWebsocketSessionManager         m_instance;
    BeastWebsocketSessions                      m_beastWebsocketSessions;
    std::mutex                                  m_beastWebsocketSessionsMutex;
    std::function<void(const std::string &msg)> m_webMsgHandler;

    BeastWebsocketSessionManager();
    void dummyWebMsgHandler(const std::string &msg);
};
