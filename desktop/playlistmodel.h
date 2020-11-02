#ifndef PLAYLISTMODEL_H
#define PLAYLISTMODEL_H

#include <QAbstractListModel>

class PlaylistModel : public QAbstractListModel
{
    Q_OBJECT
    
public:
    explicit PlaylistModel(QObject *parent = nullptr);
    
    // Basic functionality:
    int rowCount(const QModelIndex &parent = QModelIndex()) const override;
    
    QVariant data(const QModelIndex &index, int role = Qt::DisplayRole) const override;
    
    // Editable:
    bool setData(const QModelIndex &index, const QVariant &value,
                 int role = Qt::EditRole) override;
    
    Qt::ItemFlags flags(const QModelIndex& index) const override;
    
    // Add data:
    bool insertRows(int row, int count, const QModelIndex &parent = QModelIndex()) override;
    
    // Remove data:
    bool removeRows(int row, int count, const QModelIndex &parent = QModelIndex()) override;

    void setPlaylists(const QStringList &playlists);

    void addPlaylist();

    void deletePlaylist(int row);

    void filterPlaylist(const QString &keyword);

    bool isFilteredPlaylists() const;

private:
    bool        m_isFiltered {false};
    QStringList m_filteredPlaylists;
    QStringList m_playlists;
};

#endif // PLAYLISTMODEL_H
