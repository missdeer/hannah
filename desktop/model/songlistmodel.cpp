#include "songlistmodel.h"

SonglistModel::SonglistModel(QObject *parent)
    : QAbstractTableModel(parent)
{
}

int SonglistModel::rowCount(const QModelIndex &parent) const
{
    if (parent.isValid())
        return 0;

    return m_songs.length();
}

int SonglistModel::columnCount(const QModelIndex &parent) const
{
    if (parent.isValid())
        return 0;

    return 1;
}

QVariant SonglistModel::data(const QModelIndex &index, int role) const
{
    if (!index.isValid())
        return QVariant();

    if (role != Qt::DisplayRole)
        return QVariant();

    if (index.row() < 0 || index.row() >= m_songs.length())
        return QVariant();

    return QVariant(m_songs[index.row()]);
}

bool SonglistModel::setData(const QModelIndex &index, const QVariant &value, int role)
{
    if (data(index, role) != value) {
        // FIXME: Implement me!
        emit dataChanged(index, index, QVector<int>() << role);
        return true;
    }
    return false;
}

Qt::ItemFlags SonglistModel::flags(const QModelIndex &index) const
{
    if (!index.isValid())
        return Qt::NoItemFlags;

    return QAbstractItemModel::flags(index);
}

bool SonglistModel::insertRows(int row, int count, const QModelIndex &parent)
{
    beginInsertRows(parent, row, row + count - 1);
    for (int i = row; i < row + count; i++)
        m_songs.insert(i, "");
    endInsertRows();
    return true;
}

bool SonglistModel::removeRows(int row, int count, const QModelIndex &parent)
{
    beginRemoveRows(parent, row, row + count - 1);
    for (int i = row + count; i > row; i--)
        m_songs.removeAt(i - 1);
    endRemoveRows();
    return true;
}

QVariant SonglistModel::headerData(int section, Qt::Orientation orientation, int role) const
{
    if (orientation == Qt::Horizontal && role == Qt::DisplayRole)
    {
        QMap<int, QString> m = {
            {0, tr("URI")},
        };
        auto it = m.find(section);
        if (it != m.end())
            return QVariant(it.value());
    }

    return QVariant();
}

void SonglistModel::setSongList(const QStringList &songs)
{
    int rc  = m_songs.length();
    m_songs = songs;
    emit dataChanged(index(0, 0), index(rc - 1, 0));
}

void SonglistModel::appendToSonglist(const QStringList &s) {}

void SonglistModel::clearAndAddToSonglist(const QStringList &s) {}

void SonglistModel::appendToSonglistFile(const QStringList &s) {}

void SonglistModel::clearAndAddToSonglistFile(const QStringList &s) {}

void SonglistModel::deleteSongs(QList<int> rows) {}
