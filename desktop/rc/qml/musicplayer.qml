import QtQuick 2.15
import QtQuick.Layouts 1.12
import QtQuick.Controls 2.15
import QtQuick.Window 2.15

ApplicationWindow {
    id: window
    width: 1280
    height: 720
    visible: true
    title: "Hannah's Builtin Music Player"

    Component.onCompleted: {
    }

    Connections {
        target: playerCore
        function onShowPlayer() {
            window.show();
        }
    }

    Shortcut {
        sequence: "Ctrl+Q"
        onActivated: playerCore.onQuit()
    }

    header: ToolBar {
        RowLayout {
            id: headerRowLayout
            anchors.fill: parent
            spacing: 0

            ToolButton {
                icon.name: "grid"
                onClicked: playerCore.onShowPlaylists();
            }
            ToolButton {
                icon.name: "settings"
                onClicked: playerCore.onSettings();
            }
            ToolButton {
                icon.name: "filter"
                onClicked: playerCore.onFilter();
            }
            ToolButton {
                icon.name: "message"
                onClicked: playerCore.onMessage();
            }
            ToolButton {
                icon.name: "music"
                onClicked: playerCore.onMusic();
            }
            ToolButton {
                icon.name: "cloud"
                onClicked: playerCore.onCloud();
            }
            ToolButton {
                icon.name: "bluetooth"
                onClicked: playerCore.onBluetooth();
            }
            ToolButton {
                icon.name: "cart"
                onClicked: playerCore.onCart();
            }

            Item {
                Layout.fillWidth: true
            }

            ToolButton {
                icon.name: "power"
                onClicked: playerCore.onQuit();
            }
        }
    }

    Label {
        text: "Hannah - Listen music"
        font.pixelSize: Qt.application.font.pixelSize * 1.3
        anchors.centerIn: header
        z: header.z + 1
    }

    RowLayout {
        spacing: 115
        anchors.fill: parent
        anchors.margins: 70

        ColumnLayout {
            spacing: 0
            Layout.preferredWidth: 230

            RowLayout {
                Layout.maximumHeight: 170

                ColumnLayout {
                    Label {
                        text: "15 dB"
                        Layout.fillHeight: true
                    }
                    Label {
                        text: "7 dB"
                        Layout.fillHeight: true
                    }
                    Label {
                        text: "0 dB"
                        Layout.fillHeight: true
                    }
                    Label {
                        text: "-7 dB"
                        Layout.fillHeight: true
                    }
                    Label {
                        text: "-15 dB"
                        Layout.fillHeight: true
                    }
                }

                Slider {
                    id: eq0
                    from: -15
                    to: 15
                    stepSize: 1.0
                    snapMode: Slider.SnapAlways
                    value: playerCore.eq0
                    orientation: Qt.Vertical

                    Layout.fillWidth: true
                    Layout.fillHeight: true
                }
                Slider {
                    id: eq1
                    from: -15
                    to: 15
                    stepSize: 1.0
                    snapMode: Slider.SnapAlways
                    value: playerCore.eq1
                    orientation: Qt.Vertical

                    Layout.fillWidth: true
                    Layout.fillHeight: true
                }
                Slider {
                    id: eq2
                    from: -15
                    to: 15
                    stepSize: 1.0
                    snapMode: Slider.SnapAlways
                    value: playerCore.eq2
                    orientation: Qt.Vertical

                    Layout.fillWidth: true
                    Layout.fillHeight: true
                }
                Slider {
                    id: eq3
                    from: -15
                    to: 15
                    stepSize: 1.0
                    snapMode: Slider.SnapAlways
                    value: playerCore.eq3
                    orientation: Qt.Vertical

                    Layout.fillWidth: true
                    Layout.fillHeight: true
                }
                Slider {
                    id: eq4
                    from: -15
                    to: 15
                    stepSize: 1.0
                    snapMode: Slider.SnapAlways
                    value: playerCore.eq4
                    orientation: Qt.Vertical

                    Layout.fillWidth: true
                    Layout.fillHeight: true
                }
                Slider {
                    id: eq5
                    from: -15
                    to: 15
                    stepSize: 1.0
                    snapMode: Slider.SnapAlways
                    value: playerCore.eq5
                    orientation: Qt.Vertical

                    Layout.fillWidth: true
                    Layout.fillHeight: true
                }
                Slider {
                    id: eq6
                    from: -15
                    to: 15
                    stepSize: 1.0
                    snapMode: Slider.SnapAlways
                    value: playerCore.eq6
                    orientation: Qt.Vertical

                    Layout.fillWidth: true
                    Layout.fillHeight: true
                }
                Slider {
                    id: eq7
                    from: -15
                    to: 15
                    stepSize: 1.0
                    snapMode: Slider.SnapAlways
                    value: playerCore.eq7
                    orientation: Qt.Vertical

                    Layout.fillWidth: true
                    Layout.fillHeight: true
                }
                Slider {
                    id: eq8
                    from: -15
                    to: 15
                    stepSize: 1.0
                    snapMode: Slider.SnapAlways
                    value: playerCore.eq8
                    orientation: Qt.Vertical

                    Layout.fillWidth: true
                    Layout.fillHeight: true
                }
                Slider {
                    id: eq9
                    from: -15
                    to: 15
                    stepSize: 1.0
                    snapMode: Slider.SnapAlways
                    value: playerCore.eq9
                    orientation: Qt.Vertical

                    Layout.fillWidth: true
                    Layout.fillHeight: true
                }
            }

            RowLayout {
                spacing: 10
                Layout.topMargin: 23

                ComboBox {
                    currentIndex: 0
                    model: ["Default", "Pop", "Rocks", "Electronic","Classical","Metal","Dance","Country","Jazz","Bruce","Nostalgia","Opera","Voice"]
                    Layout.fillWidth: true
                    onCurrentIndexChanged: playerCore.presetEQChanged(currentIndex);
                }

                Button {
                    icon.name: "folder"
                    onClicked: playerCore.onOpenPreset();
                }

                Button {
                    icon.name: "save"
                    enabled: false
                    onClicked: playerCore.onSavePreset();
                }
            }

            Dial {
                Layout.alignment: Qt.AlignHCenter
                Layout.topMargin: 50
                value: playerCore.volumn
            }

            Label {
                text: "Volume"

                Layout.alignment: Qt.AlignHCenter
                Layout.topMargin: 12
            }
        }

        ColumnLayout {
            spacing: 26
            Layout.preferredWidth: 230

            Item {
                Layout.fillWidth: true
                Layout.fillHeight: true

                Image {
                    anchors.fill: parent
                    fillMode: Image.PreserveAspectCrop
                    source: playerCore.coverUrl
                }
            }

            Item {
                id: songLabelContainer
                clip: true

                Layout.fillWidth: true
                Layout.preferredHeight: songNameLabel.implicitHeight

                SequentialAnimation {
                    running: true
                    loops: Animation.Infinite

                    PauseAnimation {
                        duration: 2000
                    }
                    ParallelAnimation {
                        XAnimator {
                            target: songNameLabel
                            from: 0
                            to: songLabelContainer.width - songNameLabel.implicitWidth
                            duration: 5000
                        }
                        OpacityAnimator {
                            target: leftGradient
                            from: 0
                            to: 1
                        }
                    }
                    OpacityAnimator {
                        target: rightGradient
                        from: 1
                        to: 0
                    }
                    PauseAnimation {
                        duration: 1000
                    }
                    OpacityAnimator {
                        target: rightGradient
                        from: 0
                        to: 1
                    }
                    ParallelAnimation {
                        XAnimator {
                            target: songNameLabel
                            from: songLabelContainer.width - songNameLabel.implicitWidth
                            to: 0
                            duration: 5000
                        }
                        OpacityAnimator {
                            target: leftGradient
                            from: 0
                            to: 1
                        }
                    }
                    OpacityAnimator {
                        target: leftGradient
                        from: 1
                        to: 0
                    }
                }

                Rectangle {
                    id: leftGradient
                    gradient: Gradient {
                        GradientStop {
                            position: 0
                            color: "#dfe4ea"
                        }
                        GradientStop {
                            position: 1
                            color: "#00dfe4ea"
                        }
                    }

                    width: height
                    height: parent.height
                    anchors.left: parent.left
                    z: 1
                    rotation: -90
                    opacity: 0
                }

                Label {
                    id: songNameLabel
                    text: playerCore.songName
                    font.pixelSize: Qt.application.font.pixelSize * 1.4
                }

                Rectangle {
                    id: rightGradient
                    gradient: Gradient {
                        GradientStop {
                            position: 0
                            color: "#00dfe4ea"
                        }
                        GradientStop {
                            position: 1
                            color: "#dfe4ea"
                        }
                    }

                    width: height
                    height: parent.height
                    anchors.right: parent.right
                    rotation: -90
                }
            }

            RowLayout {
                spacing: 8
                Layout.alignment: Qt.AlignHCenter

                RoundButton {
                    icon.name: "favorite"
                    icon.width: 32
                    icon.height: 32
                    onClicked: playerCore.onFavorite()
                }
                RoundButton {
                    icon.name: "stop"
                    icon.width: 32
                    icon.height: 32
                    onClicked: playerCore.onStop()
                }
                RoundButton {
                    icon.name: "previous"
                    icon.width: 32
                    icon.height: 32
                    onClicked: playerCore.onPrevious()
                }
                RoundButton {
                    icon.name: "pause"
                    icon.width: 32
                    icon.height: 32
                    onClicked: playerCore.onPause()
                }
                RoundButton {
                    icon.name: "next"
                    icon.width: 32
                    icon.height: 32
                    onClicked: playerCore.onNext();
                }
                RoundButton {
                    icon.name: "repeat"
                    icon.width: 32
                    icon.height: 32
                    onClicked: playerCore.onRepeat()
                }
                RoundButton {
                    icon.name: "shuffle"
                    icon.width: 32
                    icon.height: 32
                    onClicked: playerCore.onShuffle()
                }
            }

            Slider {
                id: seekSlider
                value: playerCore.progress
                to: 261

                Layout.fillWidth: true

                ToolTip {
                    parent: seekSlider.handle
                    visible: seekSlider.pressed
                    text: pad(Math.floor(value / 60)) + ":" + pad(Math.floor(value % 60))
                    y: parent.height

                    readonly property int value: seekSlider.valueAt(seekSlider.position)

                    function pad(number) {
                        if (number <= 9)
                            return "0" + number;
                        return number;
                    }
                }
            }
        }

        ColumnLayout {
            spacing: 16
            Layout.preferredWidth: 230

            ButtonGroup {
                buttons: libraryRowLayout.children
            }

            RowLayout {
                id: libraryRowLayout
                Layout.alignment: Qt.AlignHCenter

                Button {
                    text: "Files"
                    checked: true
                    onClicked: playerCore.onSwitchFiles()
                }
                Button {
                    text: "Playlists"
                    checkable: true
                    onClicked: playerCore.onSwitchPlaylists()
                }
                Button {
                    text: "Favourites"
                    checkable: true
                    onClicked: playerCore.onSwitchFavourites()
                }
            }

            RowLayout {
                TextField {
                    Layout.fillWidth: true
                }
                Button {
                    icon.name: "folder"
                    onClicked: playerCore.onOpenFile()
                }
            }

            Frame {
                id: filesFrame
                leftPadding: 1
                rightPadding: 1

                Layout.fillWidth: true
                Layout.fillHeight: true

                ListView {
                    id: filesListView
                    clip: true
                    anchors.fill: parent
                    model: ListModel {
                        Component.onCompleted: {
                            for (var i = 0; i < 100; ++i) {
                                append({
                                   author: "Author",
                                   album: "Album",
                                   track: "Track 0" + (i % 9 + 1),
                                });
                            }
                        }
                    }
                    delegate: ItemDelegate {
                        text: model.author + " - " + model.album + " - " + model.track
                        width: filesListView.width
                    }

                    ScrollBar.vertical: ScrollBar {
                        parent: filesFrame
                        policy: ScrollBar.AlwaysOn
                        anchors.top: parent.top
                        anchors.topMargin: filesFrame.topPadding
                        anchors.right: parent.right
                        anchors.rightMargin: 1
                        anchors.bottom: parent.bottom
                        anchors.bottomMargin: filesFrame.bottomPadding
                    }
                }
            }
        }
    }
}
