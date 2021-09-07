const roomAPI = "/api/room"
const createRoomForm = document.querySelector("#create-room")

async function createRoom(e){
    e.preventDefault()
    const roomData = {
        name: this.querySelector("#room").value
    }
    const res = await fetch(roomAPI, {
        method: "POST",
        body: JSON.stringify(roomData),
        headers: {'Content-Type': 'application/json'}
    })
    const data = await res.json()
    if(data.ok){
        location = "/my/rooms"
    } else{
        this.querySelector(".msg").style.color = "red"
        this.querySelector(".msg").textContent = data.msg
    }
}

createRoomForm.addEventListener("submit", createRoom)