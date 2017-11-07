package tcpClient

import (
	"436bin/a1/model"
	"fmt"
	"testing"
)

func TestRemoveTab(t *testing.T) {
	tabName1 := "foo"
	tabName2 := "bar"
	ui := model.GetUIInstance()
	ui.ChatTabs = append(ui.ChatTabs, &model.ClientChatTab{Name: tabName1})
	ui.ChatTabs = append(ui.ChatTabs, &model.ClientChatTab{Name: tabName2})
	// Ensure set up correctly
	if len(ui.ChatTabs) != 2 {
		fmt.Println("Chat tabs not set up properly")
		t.Fail()
	}

	// Execute block
	err := ui.RemoveTab(tabName1)

	// Assert block
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	if len(ui.ChatTabs) != 1 || ui.ChatTabs[0].Name != tabName2 {
		fmt.Println("Failed to remove the desired tab by name")
		t.Fail()
	}
}

func TestRemoveTabReturnsErrorIfTabDoesNotExist(t *testing.T) {
	nonExistantTabName := "foo"
	ui := model.GetUIInstance()
	// Ensure set up correctly
	if len(ui.ChatTabs) != 0 {
		fmt.Println("Chat tabs not set up properly")
		t.Fail()
	}

	// Execute block
	err := ui.RemoveTab(nonExistantTabName)

	// Assert block
	if err == nil {
		fmt.Println("Removing non-existant tab by name did not throw expected error")
		t.Fail()
	}
	if err.Error() != fmt.Sprintf("Could not find tab named %s", nonExistantTabName) {
		fmt.Println(err.Error())
		t.Fail()
	}
}

func TestGetTabByName(t *testing.T) {
	tabName1 := "foo"
	tabName2 := "bar"
	expectedTab := &model.ClientChatTab{Name: tabName1}
	tab2 := &model.ClientChatTab{Name: tabName2}
	ui := model.GetUIInstance()
	ui.ChatTabs = append(ui.ChatTabs, expectedTab)
	ui.ChatTabs = append(ui.ChatTabs, tab2)
	// Ensure set up correctly
	if len(ui.ChatTabs) != 2 {
		fmt.Println("Chat tabs not set up properly")
		t.Fail()
	}

	// Execute block
	actualTab, err := ui.GetTabByName(tabName1)

	// Assert block
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	if actualTab != expectedTab {
		fmt.Println("Failed to retrieve the expected tab by name")
		t.Fail()
	}
}

func TestGetTabByNameReturnsErrorIfNoTabByNameExists(t *testing.T) {
	nonExistantTabName := "Randomname123"
	tabName1 := "foo"
	tabName2 := "bar"
	tab1 := &model.ClientChatTab{Name: tabName1}
	tab2 := &model.ClientChatTab{Name: tabName2}
	ui := model.GetUIInstance()
	ui.ChatTabs = append(ui.ChatTabs, tab1)
	ui.ChatTabs = append(ui.ChatTabs, tab2)
	// Ensure set up correctly
	if len(ui.ChatTabs) != 2 {
		fmt.Println("Chat tabs not set up properly")
		t.Fail()
	}

	// Execute block
	actualTab, err := ui.GetTabByName(nonExistantTabName)

	// Assert block
	if actualTab != nil {
		fmt.Println("Found unexpected tab for non-existant name")
		t.Fail()
	}
	if err == nil {
		fmt.Println("Error was not thrown when tab was not found")
		t.Fail()
	}
	if err.Error() != fmt.Sprintf("Could not find tab named %s", nonExistantTabName) {
		fmt.Println(err.Error())
		t.Fail()
	}
}
