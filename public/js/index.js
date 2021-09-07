const roomAPI = "/api/rooms"
page = 1
let firstFetch = true

const roomLst = document.querySelector("#room-list")
const pageBar = document.querySelector("#page-bar")
const prevBtn = pageBar.querySelector("#prev-page")
const nextBtn = pageBar.querySelector("#next-page")
const pageGroup = pageBar.querySelector("#page-group")

async function getRooms(){
    const res = await fetch(`${roomAPI}?page=${page}`)
    const data = await res.json()
    if(data.error){
        console.log(data.msg)
    }else{
        roomLst.innerHTML = `
        <ul class="list-group list-group-horizontal">
            <li class="list-group-item col-4">聊天室名稱</li>
            <li class="list-group-item col-4">室長</li>
            <li class="list-group-item col-4">操作</li>
        </ul>
        `

        data.rooms.forEach(room => {
            // console.log(room)
            
            // 包一筆room資料 
            const roomContainer = document.createElement("ul")
            roomContainer.className = "list-group list-group-horizontal"
            
            // 房間名稱
            const roomNameContainer = createContainer()
            const roomName = document.createElement("a")
            roomName.href = room.url
            roomName.innerText = room.name
            roomName.target = "_blank"
            roomNameContainer.append(roomName)

            // 室長
            const owner = createContainer()
            owner.innerText = room.owner

            // 房間入口
            const entryContainer = createContainer()
            const entry = document.createElement("a")
            entry.href = room.url
            entry.innerText = "進入"
            entry.target = "_blank"
            entryContainer.append(entry)

            roomContainer.append(roomNameContainer, owner, entryContainer)
            roomLst.append(roomContainer)
        })
        
        if(page == 1){
            prevBtn.disabled = true
            prevBtn.classList.remove("btn-primary")
            prevBtn.classList.add("btn-secondary")
        } else{
            prevBtn.disabled = false
            prevBtn.classList.remove("btn-secondary")
            prevBtn.classList.add("btn-primary")
        }

        if(page == data.pageStatus.maxPage){
            nextBtn.disabled = true
            nextBtn.classList.remove("btn-primary")
            nextBtn.classList.add("btn-secondary")
        } else{
            nextBtn.disabled = false
            nextBtn.classList.remove("btn-secondary")
            nextBtn.classList.add("btn-primary")
        }

        if(firstFetch){
            const pageGroupUl = pageGroup.querySelector("ul")
            pageGroupUl.style.minWidth = "5rem"
            for(let i = 1; i <= data.pageStatus.maxPage; i++){
                const pageContainer = document.createElement("li")
                const pageBtn = document.createElement("button")
                pageBtn.className = "btn d-block mx-auto"
                pageBtn.innerText = `第${i}頁`
                pageContainer.append(pageBtn)
                pageGroupUl.append(pageContainer)

                pageBtn.addEventListener("click", ()=>{
                    page = i
                    getRooms()
                })
            }
            firstFetch = false
        }
    }
}

function createContainer(){
    const li = document.createElement("li")
    li.className = "list-group-item col-4"
    return li
}

getRooms()

prevBtn.addEventListener("click", () => {
    page = page-1
    getRooms()
})

nextBtn.addEventListener("click", () => {
    page = page+1
    getRooms()
})
