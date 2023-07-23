#pragma once

#include <map>
#include <string>
#include <vector>

using QueryPairType = std::map<std::string, std::string>;

//! UrlQuery parse something like https://www.example.com/some/path?id=xxx&name=yyy&age=zzz
class UrlQuery
{
public:
    explicit UrlQuery(const std::string &uri);

    QueryPairType::iterator begin();
    QueryPairType::iterator end();
    bool                    contains(const std::string &key);
    const std::string      &value(const std::string &key);
    int                     intValue(const std::string &key);
    bool                    isInt(const std::string &key);
    std::int64_t            int64Value(const std::string &key);
    bool                    isInt64(const std::string &key);
    bool                    boolValue(const std::string &key);
    bool                    isBool(const std::string &key);

private:
    QueryPairType m_queries;

    void parse(const std::string &uri);
};

//! UrlParam represent something like https://www.exmample.com/some/{id}/path
class UrlParam
{
public:
    void               push_back(const std::string &value);
    bool               isInt(size_t index);
    int                intAt(size_t index);
    bool               isInt64(size_t index);
    std::int64_t       int64At(size_t index);
    const std::string &at(size_t index);
    size_t             size() const;
    bool               empty() const;

private:
    std::vector<std::string> m_values;
};