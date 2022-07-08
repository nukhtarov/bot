package main

var admins = map[int]int{
	237286647: 237286647,
	162667568: 162667568,
}

func isAdmin(id int) bool {
	_, has := admins[id]
	return has

}
