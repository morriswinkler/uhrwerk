import QtQuick 1.0

Rectangle {
	id: page
	width: 640
	height: 400
	//width: 1280 // chalk elec screen
	//height: 800 // chalk elec screen

	SimpleViewManager {
		id: viewManager 
		anchors.fill: parent

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

				Rectangle {
					id: loginButton
					anchors.top: parent.top
					anchors.bottom: parent.bottom
					anchors.right: parent.right
					anchors.topMargin: 10
					anchors.bottomMargin: 10
					anchors.rightMargin: 10
					width: 200
					color: "red"

					Text {
						id: loginButtonText
						text: 
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

			MouseArea {
				anchors.fill: parent
				onClicked: view2.visible = true
	    	}
		}

		Rectangle {
			id: view2
			width: parent.width
			height: parent.height
			color: "#ff00ffff"

			Text {
				id: view1Text
				text: "Go to main view"
				color: "white"
				anchors.centerIn: parent
	    	}
	    		
	    	MouseArea {
				anchors.fill: parent
				onClicked: view1.visible = true
	    	}
		}
	}	
}