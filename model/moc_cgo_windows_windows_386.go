package model

/*
#cgo CFLAGS: -pipe -fno-keep-inline-dllexport -O2 -Wall -Wextra -DUNICODE -DQT_NO_DEBUG -DQT_MULTIMEDIAWIDGETS_LIB -DQT_MULTIMEDIA_LIB -DQT_QUICKWIDGETS_LIB -DQT_WIDGETS_LIB -DQT_QUICK_LIB -DQT_GUI_LIB -DQT_QML_LIB -DQT_NETWORK_LIB -DQT_DBUS_LIB -DQT_TESTLIB_LIB -DQT_CORE_LIB 
#cgo CXXFLAGS: -pipe -fno-keep-inline-dllexport -O2 -std=gnu++11 -frtti -Wall -Wextra -fexceptions -mthreads -DUNICODE -DQT_NO_DEBUG -DQT_MULTIMEDIAWIDGETS_LIB -DQT_MULTIMEDIA_LIB -DQT_QUICKWIDGETS_LIB -DQT_WIDGETS_LIB -DQT_QUICK_LIB -DQT_GUI_LIB -DQT_QML_LIB -DQT_NETWORK_LIB -DQT_DBUS_LIB -DQT_TESTLIB_LIB -DQT_CORE_LIB 
#cgo CXXFLAGS: -I../../golang-chat -I. -I/usr/lib/mxe/usr/i686-w64-mingw32.shared/qt5/include -I/usr/lib/mxe/usr/i686-w64-mingw32.shared/qt5/include/QtMultimediaWidgets -I/usr/lib/mxe/usr/i686-w64-mingw32.shared/qt5/include/QtMultimedia -I/usr/lib/mxe/usr/i686-w64-mingw32.shared/qt5/include/QtQuickWidgets -I/usr/lib/mxe/usr/i686-w64-mingw32.shared/qt5/include/QtWidgets -I/usr/lib/mxe/usr/i686-w64-mingw32.shared/qt5/include/QtQuick -I/usr/lib/mxe/usr/i686-w64-mingw32.shared/qt5/include/QtGui -I/usr/lib/mxe/usr/i686-w64-mingw32.shared/qt5/include/QtQml -I/usr/lib/mxe/usr/i686-w64-mingw32.shared/qt5/include/QtNetwork -I/usr/lib/mxe/usr/i686-w64-mingw32.shared/qt5/include/QtDBus -I/usr/lib/mxe/usr/i686-w64-mingw32.shared/qt5/include/QtTest -I/usr/lib/mxe/usr/i686-w64-mingw32.shared/qt5/include/QtCore -Irelease -I/usr/lib/mxe/usr/i686-w64-mingw32.shared/qt5/mkspecs/win32-g++
#cgo LDFLAGS:        -Wl,-s -Wl,-subsystem,console -mthreads
#cgo LDFLAGS:        -lopengl32 -L/usr/lib/mxe/usr/i686-w64-mingw32.shared/qt5/lib -lQt5MultimediaWidgets -lQt5Multimedia -lQt5QuickWidgets -lQt5Widgets -lQt5Quick -lQt5Gui -lQt5Qml -lQt5Network -lQt5DBus -lQt5Test -lQt5Core
#cgo LDFLAGS: -Wl,--allow-multiple-definition
*/
import "C"
