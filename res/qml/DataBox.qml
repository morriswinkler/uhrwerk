import QtQuick 1.0

Item {
	id: container
	property alias title: title.text
	property alias bookedTillKey: bookedTillKey.text
	property alias bookedTillValue: bookedTillValue.text
	width: parent.width / parent.columns
	height: parent.height / parent.rows

	//FontLoader { id: localFont; source: "fonts/tarzeau_ocr_a.ttf" }
	FontLoader { id: localFont; name: "Helvetica" }

	Rectangle {
		id: rect
		anchors.fill: parent
		gradient: Gradient {
			GradientStop { position: 0.0; color: "#ffffffff" }
			GradientStop { position: 1.0; color: "#ff999999" }
		}
	}
	Rectangle {
		anchors.fill: parent
		anchors.topMargin: 10
		anchors.rightMargin: 10
		anchors.bottomMargin: 10
		anchors.leftMargin: 10
		color: "#00000000"
	
		Text {
			id: title
			color: "#ffC45757"
			width: parent.width
			elide: Text.ElideLeft
			font { family: localFont.name; pointSize: 19; capitalization: Font.Capitalize; bold: true }
		}
	
		Row {
			id: bookedTillRow
			anchors.top: title.bottom
			anchors.right: parent.right
			anchors.left: parent.left
	
			Text {
				id: bookedTillKey
				x: 2
				y: 4
				anchors.left: parent.left
				width: bookedTillKey.paintedWidth
				text: "..."
				color: "#ff000000"
				font { family: localFont.name; pointSize: 13; capitalization: Font.Capitalize; bold: false }
			}
	
			Image {
				source: "images/dotted-line-1-2.png"
				fillMode: Image.TileHorizontally
				height: 1;
				anchors.left: bookedTillKey.right
				anchors.right: bookedTillValueContainer.left
				anchors.bottom: bookedTillKey.bottom
				anchors.leftMargin: 4
				anchors.rightMargin: 4
				anchors.bottomMargin: 2
			}
	
			Rectangle {
				id: bookedTillValueContainer
				anchors.right: bookedTillRow.right
				width: bookedTillValue.paintedWidth + 8
				height: bookedTillValue.paintedHeight + 6
				color: "#ffC45757"
				radius: 5
	
				Text {
					x: 4
					y: 4
					clip: true
					id: bookedTillValue
					text: "yes..."
					color: "#ffffffff"
					font { family: localFont.name; pointSize: 13; capitalization: Font.Capitalize; bold: false }
				}
			}
		}
	}

	//Component.onCompleted: console.log("index: " + index)
}