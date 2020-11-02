#include "playlistmodel.h"

PlaylistModel::PlaylistModel(QObject *parent)
    : QAbstractListModel(parent)
{
}

int PlaylistModel::rowCount(const QModelIndex &parent) const
{
    // For list models only the root node (an invalid parent) should return the list's size. For all
    // other (valid) parents, rowCount() should return 0 so that it does not become a tree model.
    if (parent.isValid())
        return 0;

    return m_isFiltered ? m_filteredPlaylists.length() : m_playlists.length();
}

QVariant PlaylistModel::data(const QModelIndex &index, int role) const
{
    if (!index.isValid())
        return QVariant();

    if (role != Qt::DisplayRole)
        return QVariant();

    if (index.row() < 0 || index.row() >= rowCount())
        return QVariant();

    return QVariant(m_isFiltered ? m_filteredPlaylists[index.row()] : m_playlists[index.row()]);
}

bool PlaylistModel::setData(const QModelIndex &index, const QVariant &value, int role)
{
    if (data(index, role) != value) {
        // FIXME: Implement me!
        emit dataChanged(index, index, QVector<int>() << role);
        return true;
    }
    return false;
}

Qt::ItemFlags PlaylistModel::flags(const QModelIndex &index) const
{
    if (!index.isValid())
        return Qt::NoItemFlags;

    return m_isFiltered ? QAbstractListModel::flags(index) : Qt::ItemIsEditable; // FIXME: Implement me!
}

bool PlaylistModel::insertRows(int row, int count, const QModelIndex &parent)
{
    beginInsertRows(parent, row, row + count - 1);
    for (int i = row; i < row + count; i++)
        m_playlists.insert(i, "");
    endInsertRows();
    return true;
}

bool PlaylistModel::removeRows(int row, int count, const QModelIndex &parent)
{
    beginRemoveRows(parent, row, row + count - 1);
    for (int i = row + count; i > row; i--)
        m_playlists.removeAt(i - 1);
    endRemoveRows();
    return true;
}

void PlaylistModel::setPlaylists(const QStringList &playlists)
{
    int rc       = m_playlists.length();
    m_playlists  = playlists;
    m_isFiltered = false;
    emit dataChanged(index(0, 0), index(rc - 1, 0));
}

void PlaylistModel::addPlaylist()
{
    int rc = rowCount();
    insertRows(rc, 1);
}

void PlaylistModel::deletePlaylist(int row)
{
    QString playlist = m_playlists[row];
    m_playlists.removeAt(row);
    removeRows(row, 1);
}

void PlaylistModel::filterPlaylist(const QString &keyword)
{
    int rc = rowCount();
    m_filteredPlaylists.clear();
    if (keyword.isEmpty())
    {
        m_isFiltered = false;
    }
    else
    {
        m_isFiltered = true;
        for (const auto &s : m_playlists)
        {
            if (s.contains(keyword, Qt::CaseInsensitive))
            {
                m_filteredPlaylists.append(s);
            }
        }
    }
    emit dataChanged(index(0, 0), index(rc - 1, 0));
}

bool PlaylistModel::isFilteredPlaylists() const
{
    return m_isFiltered;
}
