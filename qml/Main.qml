import QtQuick 1.0

Rectangle {
	id: page
	width: 640
	height: 400
	//width: 1280 // chalk elec screen
	//height: 800 // chalk elec screen

	// Load fonts
	FontLoader { 
		id: mainFont; 
		source: "fonts/Droid_Serif/DroidSerif.ttf" 
	}

	SimpleViewManager {
		id: viewManager 
		anchors.fill: parent

		// Main view
		Rectangle {
			id: view1
			width: parent.width
			height: parent.height

			gradient: Gradient {
				GradientStop { position: 0.0; color: "#ff222222" }
				GradientStop { position: 1.0; color: "#ff000000" }
			}
		
			Rectangle {
				id: topBar
				anchors.left: parent.left
				anchors.right: parent.right
				height: 50
				gradient: Gradient {
					GradientStop { position: 0.0; color: "#ff333333" }
					GradientStop { position: 1.0; color: "#ff222222" }
				}

				Text {
					id: demoMessage
					text: "Click on the Login button to test animated transition"
					color: "#ff999999"
					font.family: mainFont.name
					font.pointSize: 14
					anchors.verticalCenter: parent.verticalCenter
					anchors.left: parent.left
					anchors.leftMargin: 10
				}

				Rectangle {
					id: loginButton
					radius: 5
					anchors.top: parent.top
					anchors.bottom: parent.bottom
					anchors.right: parent.right
					anchors.topMargin: 10
					anchors.bottomMargin: 10
					anchors.rightMargin: 10
					width: loginButtonText.paintedWidth + 30
					color: "#ffffffff"

					Text {
						id: loginButtonText
						text: "Login"
						font.family: mainFont.name
						font.bold: false
						color: "#ff000000"
						anchors.centerIn: parent
						font.pointSize: 14
					}

					MouseArea {
						anchors.fill: parent
						onClicked: view2.visible = true
					}
				}
			}
		
			Grid {
				id: imgGrid
				anchors.top: topBar.bottom
				anchors.bottom: bottomBar.top
				anchors.left: parent.left
				anchors.right: parent.right
				rows: 2
				columns: 4
				spacing: 1
		
				Repeater {
					model: 6
					DataBox { title: "Bookers"; bookedTillKey: "Booked till:"; bookedTillValue: "10:00" }
				}
			}
		
			Rectangle {
				id: bottomBar
				anchors.left: parent.left
				anchors.right: parent.right
				anchors.bottom: parent.bottom
				height: 30
				gradient: Gradient {
					GradientStop { position: 0.0; color: "#ff333333" }
					GradientStop { position: 1.0; color: "#ff222222" }
				}
			}
		}

		Rectangle {
			id: bookingView
			width: parent.width
			height: parent.height

			
		}

		// Temp view
		Rectangle {
			id: view2
			width: parent.width
			height: parent.height
			color: "#ff00ffff"

			Text {
				id: view1Text
				text: "Click to go back to the main view"
				color: "black"
				font.family: mainFont.name
				font.pointSize: 24
				font.bold: false
				anchors.centerIn: parent
	    	}
	    		
	    	MouseArea {
				anchors.fill: parent
				onClicked: view1.visible = true
	    	}
		}
	}	
}