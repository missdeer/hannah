

#include <vector>
#include <boost/algorithm/string.hpp>

#include "UrlQuery.h"


UrlQuery::UrlQuery(const std::string &uri)
{
    parse(uri);
}

QueryPairType::iterator UrlQuery::begin()
{
    return m_queries.begin();
}

QueryPairType::iterator UrlQuery::end()
{
    return m_queries.end();
}

bool UrlQuery::contains(const std::string &key)
{
    return m_queries.find(key) != m_queries.end();
}

const std::string &UrlQuery::value(const std::string &key)
{
    return m_queries[key];
}

int UrlQuery::intValue(const std::string &key)
{
    auto iter = m_queries.find(key);
    return std::stoi(iter->second);
}

bool UrlQuery::isInt(const std::string &key)
{
    auto iter = m_queries.find(key);
    if (m_queries.end() != iter)
    {
        try
        {
            std::stoi(iter->second);
            return true;
        }
        catch (...)
        {
        }
    }
    return false;
}

std::int64_t UrlQuery::int64Value(const std::string &key)
{
    auto iter = m_queries.find(key);
    return std::stoll(iter->second);
}

bool UrlQuery::isInt64(const std::string &key)
{
    auto iter = m_queries.find(key);
    if (m_queries.end() != iter)
    {
        try
        {
            std::stoll(iter->second);
            return true;
        }
        catch (...)
        {
        }
    }
    return false;
}

bool UrlQuery::boolValue(const std::string &key)
{
    auto iter = m_queries.find(key);
    if (m_queries.end() != iter)
    {
        return boost::algorithm::iequals(iter->second, "true");
    }
    return false;
}

bool UrlQuery::isBool(const std::string &key)
{
    auto iter = m_queries.find(key);
    if (m_queries.end() != iter)
    {
        return boost::algorithm::iequals(iter->second, "true") || boost::algorithm::iequals(iter->second, "false");
    }
    return false;
}

void UrlQuery::parse(const std::string &uri)
{
    auto pos = uri.find('?');
    if (pos == std::string::npos)
    {
        return;
    }
    auto parameters = uri.substr(pos + 1);

    // split by &
    std::vector<std::string> expressions;
    boost::split(expressions, parameters, boost::is_any_of("&"));

    for (const auto &expression : expressions)
    {
        // split by =
        std::vector<std::string> kv;
        boost::split(kv, expression, boost::is_any_of("="));
        if (kv.size() == 2)
        {
            m_queries.insert(std::make_pair(kv[0], kv[1]));
        }
    }
}

void UrlParam::push_back(const std::string &value)
{
    m_values.push_back(value);
}

bool UrlParam::isInt(size_t index)
{
    if (index < m_values.size())
    {
        try
        {
            auto val = std::stoi(m_values.at(index));
            return true;
        }
        catch (...)
        {
        }
    }
    return false;
}

int UrlParam::intAt(size_t index)
{
    return std::stoi(m_values.at(index));
}

bool UrlParam::isInt64(size_t index)
{
    if (index < m_values.size())
    {
        try
        {
            auto val = std::stoll(m_values.at(index));
            return true;
        }
        catch (...)
        {
        }
    }
    return false;
}

std::int64_t UrlParam::int64At(size_t index)
{
    return std::stoll(m_values.at(index));
}

const std::string &UrlParam::at(size_t index)
{
    return m_values.at(index);
}

size_t UrlParam::size() const
{
    return m_values.size();
}

bool UrlParam::empty() const
{
    return m_values.empty();
}
