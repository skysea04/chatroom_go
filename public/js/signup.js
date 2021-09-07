// api
const userAPI = "/user"

const signupForm = document.querySelector("#signup") 
const signinForm = document.querySelector("#signin")

async function signup(e){
    e.preventDefault()
    const signupData = {
        email: this.querySelector("#signup-email").value,
        pwd: this.querySelector("#signup-pwd").value,
        rePwd: this.querySelector("#signup-re-pwd").value,
        name: this.querySelector("#signup-name").value
    }
    const res = await fetch(userAPI, {
        method: 'POST',
        body: JSON.stringify(signupData),
        headers: {'Content-Type': 'application/json'}
    })
    const data = await res.json()
    if (data.ok) {
        this.querySelector(".msg").style.color = "black"
        this.querySelector(".msg").textContent = data.msg
    } else {
        this.querySelector(".msg").style.color = "red"
        this.querySelector(".msg").textContent = data.msg
    }
}

async function signin(e){
    e.preventDefault()
    const signinData = {
        email: this.querySelector("#signin-email").value,
        pwd: this.querySelector("#signin-pwd").value,
    }
    const res = await fetch(userAPI, {
        method: "PATCH",
        body: JSON.stringify(signinData),
        headers : {'Content-Type': 'application/json'}
    })
    const data = await res.json()
    if(data.ok){
        location = "/"
    }else{
        this.querySelector(".msg").style.color = "red"
        this.querySelector(".msg").textContent = data.msg
    }
}

signupForm.addEventListener("submit", signup)
signinForm.addEventListener("submit", signin)