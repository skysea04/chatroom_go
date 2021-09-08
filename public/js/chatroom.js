// 房間名稱與室長
const roomName = document.querySelector("#room-name")
const roomOwner = document.querySelector("#room-owner")

// 聊天室資訊
const memLst = document.querySelector("#member-list")
const msgContainer = document.querySelector("#msg-container")
const msgForm = document.querySelector("#msg-form")

const roomID = parseInt(location.pathname.split("/")[2])
const wsURL = "ws://" + document.location.host + "/ws/" + roomID
const socket = new WebSocket(wsURL)

const wsAction = {
    showMembers: 1,
    join: 2,
    sendMsg: 3,
    leave: 4
}

socket.onmessage = (msg) => {
    const data = JSON.parse(msg.data)
    console.log(data)
    switch (data.action){
        case wsAction.showMembers:
            if(Array.isArray(data.data)){
                data.data.forEach(mem => {
                    addMember(mem)
                })
            }
            break

        case wsAction.join:
            addMember(data.name)
            break

        case wsAction.sendMsg:
            const msg = data.name + "：" + data.msg
            receiveMsg(msg)
            break

        case wsAction.leave:
            const myname = document.querySelector(`#user-${data.name}`)
            memLst.removeChild(myname)
    }
}

function receiveMsg(msg){
    const msgField = document.createElement("li")
    msgField.className = "list-group-item border-0"
    msgField.textContent = msg
    msgContainer.prepend(msgField)
}

function sendMsg(e){
    e.preventDefault()
    const msgField = this.querySelector("textarea")
    if(msgField.value === "") return
    const msgData = {
        action: wsAction.sendMsg,
        msg: msgField.value
    }
    msgField.value = ""
    console.log(msgData)
    socket.send(JSON.stringify(msgData))
}

function addMember(name){
    const member = document.createElement("li")
    member.className = "list-group-item"
    member.id = `user-${name}`
    member.textContent = name
    memLst.append(member)
}

async function getRoomInfo(){
    const res = await fetch(`/api/room/${roomID}`)
    const data = await res.json()
    if(data.ok){
        roomName.textContent = data.name
        roomOwner.textContent = "室長：" + data.owner
    } else{
        roomName.textContent = data.msg
        roomOwner.textContent = data.msg
    }
    
} 

msgForm.addEventListener("submit", sendMsg)

getRoomInfo()