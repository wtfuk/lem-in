package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Room struct {
	Name        string
	X           int
	Y           int
	Connections []*Room
	Visited     bool
	Distance    int
}

type AntHill struct {
	Rooms     []*Room
	StartRoom *Room
	EndRoom   *Room
	Ants      int
}

var ah AntHill

func ReadFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var originalFileLines []string
	for scanner.Scan() {
		originalFileLines = append(originalFileLines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return originalFileLines, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . <filename>")
		return
	}
	originalFileLines, err := ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	// Remove the comments from the original file lines
	filteredLines := RemoveComments(originalFileLines)

	// check length of slice to be minimum 6: 1st line is number of ants, 2nd  and 3rd line is start room, 4th and 5th line is end room, 6th line is a link
	if len(filteredLines) < 6 {
		NoGo("")
	}

	// check if first line is a number
	if !IsNumber(filteredLines[0]) {
		NoGo("")
	}

	// convert first line to int and store in AntNum
	ah.Ants, _ = strconv.Atoi(filteredLines[0])
	filteredLines = filteredLines[1:]

	// check if number of ants is valid
	if ah.Ants <= 0 {
		NoGo("")
	}

	No2Dashes(filteredLines)
	No3Spaces(filteredLines)
	NoDuplicateLines(filteredLines)
	NoHashInLastLine(filteredLines)

	// extract start room
	ExtractStartRoom(filteredLines)
	filteredLines = DeleteStartRoom(filteredLines)

	// extract end room
	ExtractEndRoom(filteredLines)
	filteredLines = DeleteEndRoom(filteredLines)

	// extract rooms
	ExtractRooms(filteredLines)
	OnlyConnections := DeleteAllRooms(filteredLines)

	// check if any room is there in the connections that is not in the rooms
	CheckRoomsInConnectionsPresent(OnlyConnections, GetAllRoomNames(&ah))

	// Print the contents of the slice with a new line after each element
	fmt.Println(strings.Join(originalFileLines, "\n") + "\n")

	// Add Connections to the rooms where a connection is in the format "room1-room2" and room1 and room2 are in the rooms
	AddConnections(OnlyConnections)

	PrintAnthill()

}

func AddConnections(OnlyConnections []string) {
	for _, connection := range OnlyConnections {
		room1Name := strings.Split(connection, "-")[0]
		room2Name := strings.Split(connection, "-")[1]
		room1 := GetRoomByName(room1Name)
		room2 := GetRoomByName(room2Name)
		room1.Connections = append(room1.Connections, room2)
		room2.Connections = append(room2.Connections, room1)
	}
}

func GetRoomByName(name string) *Room {
	for _, room := range ah.Rooms {
		if room.Name == name {
			return room
		}
	}
	return nil
}

// RemoveComments removes the comments from the original file lines
func RemoveComments(originalFileLines []string) []string {
	var filteredLines []string
	for _, line := range originalFileLines {
		if strings.HasPrefix(line, "#") && line != "##end" && line != "##start" {
			continue
		}
		filteredLines = append(filteredLines, line)
	}
	return filteredLines
}

// IsNumber checks if a string is a number
func IsNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// No2Dashes checks if there are 2 or more dashes in a line
func No2Dashes(s []string) {
	for _, line := range s {
		if len(strings.Split(line, "-")) > 2 {
			NoGo("2 or more dashes in a line are not allowed")
		}
	}
}

// No3Spaces checks if there are 3 or more spaces in a line
func No3Spaces(s []string) {
	for _, line := range s {
		if len(strings.Split(line, " ")) > 3 {
			NoGo("3 or more spaces in a line are not allowed")
		}
	}
}

// ExtractStartRoom extracts the start room from the slice
func ExtractStartRoom(s []string) {
	for i, line := range s {
		if line == "##start" {
			if i+1 < len(s) && IsRoom(s[i+1]) {
				ah.StartRoom = ConvertToRoom(s[i+1])
			} else {
				NoGo("")
			}
		}
	}
}

// ExtractEndRoom extracts the end room from the slice
func ExtractEndRoom(s []string) {
	for i, line := range s {
		if line == "##end" {
			if i+1 < len(s) && IsRoom(s[i+1]) {
				ah.EndRoom = ConvertToRoom(s[i+1])
			} else {
				NoGo("")
			}
		}
	}
}

// DeleteStartRoom deletes the start room from the slice
func DeleteStartRoom(s []string) []string {
	var filteredLines []string
	startRoomIndex := -1
	for i, line := range s {
		if i == startRoomIndex {
			continue
		}
		if line == "##start" {
			startRoomIndex = i + 1
			continue
		}
		filteredLines = append(filteredLines, line)
	}
	return filteredLines
}

// DeleteEndRoom deletes the end room from the slice
func DeleteEndRoom(s []string) []string {
	var filteredLines []string
	endRoomIndex := -1
	for i, line := range s {
		if i == endRoomIndex {
			continue
		}
		if line == "##end" {
			endRoomIndex = i + 1
			continue
		}
		filteredLines = append(filteredLines, line)
	}
	return filteredLines
}

// DeleteAllRooms deletes all the rooms from the slice
func DeleteAllRooms(s []string) []string {
	var filteredLines []string
	for _, line := range s {
		if !IsRoom(line) {
			filteredLines = append(filteredLines, line)
		}
	}
	return filteredLines
}

// NoDuplicateLines checks if there are duplicate lines in the slice
func NoDuplicateLines(s []string) {
	for i, line := range s {
		for j, line2 := range s {
			if i != j && line == line2 {
				NoGo("Duplicate lines are not allowed")
			}
		}
	}
}

// NoDuplicateCoords checks if there are duplicate coordinates in the slice
func NoDuplicateCoords(s []*Room) {
	for i, room := range s {
		for j, room2 := range s {
			if i != j && room.X == room2.X && room.Y == room2.Y {
				NoGo("Duplicate coordinates are not allowed")
			}
		}
	}
}

// ExtractRooms extracts all the rooms from the slice
func ExtractRooms(s []string) {
	var rooms []*Room
	for _, line := range s {
		if IsRoom(line) {
			rooms = append(rooms, ConvertToRoom(line))
		}
	}
	rooms = append(rooms, ah.StartRoom)
	rooms = append(rooms, ah.EndRoom)
	NoDuplicateCoords(rooms)
	// fill AntHill with the rooms
	ah.Rooms = rooms
}

// ConvertToRoom converts a string to a room
func ConvertToRoom(roomStr string) *Room {
	// split the room line into a slice
	roomStrSlice := strings.Split(roomStr, " ")
	// convert the coordinates to ints
	rName := roomStrSlice[0]
	x, _ := strconv.Atoi(roomStrSlice[1])
	y, _ := strconv.Atoi(roomStrSlice[2])
	return &Room{
		Name: rName,
		X:    x,
		Y:    y,
	}
}

// No # in last line, or it is a start or end room
func NoHashInLastLine(s []string) {
	if strings.HasPrefix(s[len(s)-1], "#") {
		NoGo("")
	}
}

// IsRoom checks if a string is a room
func IsRoom(s string) bool {
	return !((len(strings.Split(s, " ")) != 3) || !IsNumber(strings.Split(s, " ")[1]) || !IsNumber(strings.Split(s, " ")[2]))
}

func NoGo(msg string) {
	fmt.Println("ERROR: invalid data format")
	if msg != "" {
		fmt.Println(msg)
	}
	os.Exit(1)
}

func CheckRoomsInConnectionsPresent(OnlyConnections []string, AllRooms []string) {
	for _, connectionStr := range OnlyConnections {
		// split the connectionStr line into a slice of roomsnames by "-"
		roomNames := strings.Split(connectionStr, "-")
		if !Contains(AllRooms, roomNames[0]) || !Contains(AllRooms, roomNames[1]) {
			NoGo("ERROR: room in connection not present in rooms")
		}
	}
}

// Contains checks if a string is in a slice
func Contains(slice []string, elem string) bool {
	return strings.Contains(strings.Join(slice, "😎"), elem)
}

// CheckRoomsInConnectionsPresent checks if all the rooms in the connections are present in the rooms
func GetAllRoomNames(ah *AntHill) []string {
	var roomNames []string
	for _, room := range ah.Rooms {
		roomNames = append(roomNames, room.Name)
	}
	return roomNames
}

func PrintAnthill() {
	fmt.Println("Ant Hill:")
	fmt.Println("Start Room:", ah.StartRoom.Name)
	fmt.Println("End Room:", ah.EndRoom.Name)
	fmt.Println("Number of Ants:", ah.Ants)
	fmt.Println("Rooms:")
	for _, room := range ah.Rooms {
		fmt.Printf("Name: %s, Coordinates: (%d,%d), Distance: %d, Connections: [", room.Name, room.X, room.Y, room.Distance)
		for i, connection := range room.Connections {
			fmt.Printf("%s", connection.Name)
			if i < len(room.Connections)-1 {
				fmt.Printf(", ")
			}
		}
		fmt.Println("]")
	}
}

// RoomInConnections checks if a given room is in a slice of rooms
func RoomInConnections(connections []*Room, room *Room) bool {
	for _, r := range connections {
		if r == room {
			return true
		}
	}
	return false
}
