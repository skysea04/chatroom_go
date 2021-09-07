const roomAPI = "/api/my/rooms"
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
        if (data.pageStatus.maxPage != 0) {
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
        }
        
        // 第一次render頁面
        if(firstFetch){
            if(data.pageStatus.maxPage > 1){
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
            }else{
                const pageBtn = pageGroup.querySelector("button")
                disableBtn(pageBtn)
            }
            firstFetch = false
        }
        
        
        if(page == 1){
            disableBtn(prevBtn)
        } else{
            ableBtn(prevBtn)
        }

        if(page >= data.pageStatus.maxPage){
            disableBtn(nextBtn)
        } else{
            ableBtn(nextBtn)
        }

    }
}

function createContainer(){
    const li = document.createElement("li")
    li.className = "list-group-item col-4"
    return li
}

function disableBtn(btn){
    btn.disabled = true
    btn.classList.remove("btn-primary")
    btn.classList.add("btn-secondary")
}

function ableBtn(btn){
    btn.disabled = false
    btn.classList.remove("btn-secondary")
    btn.classList.add("btn-primary")
}


prevBtn.addEventListener("click", () => {
    page = page-1
    getRooms()
})

nextBtn.addEventListener("click", () => {
    page = page+1
    getRooms()
})

// 執行第一次
getRooms()