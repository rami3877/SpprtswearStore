fetch("/AllContainerAndKind", {
    method: "GET"
}).then(res => res.json()).then(data => {
    var listallContainer = Object.getOwnPropertyNames(data);
    listallContainer.forEach(element => {
        var li = document.createElement("li")
        var a = document.createElement("a")
        a.innerText = element
        a.setAttribute("href", "/" + element)
        li.appendChild(a)
        document.getElementById("navlist").appendChild(li)

        li = document.createElement("li")
        a = document.createElement("a")
        a.innerText = element
        a.setAttribute("href", "/" + element)
        li.appendChild(a)
        document.getElementById("main-contact").appendChild(li)
    });
})


function logoutUser() {

    fetch("/user/logout", {
        method: "GET"
    }).then(re => re.ok).then(d => {
            window.location.reload()
    })

}


function loginIndex() {
    const username = document.getElementById("usernameLogin").value
    const password = document.getElementById("passwordLogin").value
    fetch("/user/login", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: `{"username":"${username}","password":"${password}"}`,
    }).then(Response => Response.json()).then(data => {
        if ( data != "ok" ) {
            document.getElementById("spanErrorLogin").innerText = data
        } else {
            window.location.reload()
        }
    })

}

function ClearErrorSpan() {
    document.getElementById("spanErrorLogin").innerText = ""
    document.getElementById("spanErrorResister").innerText = ""
}


function Resister() {
    const username = document.getElementById("usernameResister").value
    const email = document.getElementById("emailResister").value
    const password = document.getElementById("passwordResister").value
    fetch("/user/register", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: `{"username":"${username}","password":"${password}","email":"${email}"}`,
    }).then(Response => Response.text()).then(data => {
        if (data.length != 0 && data != "create" ) {
            document.getElementById("spanErrorResister").innerText = data
        } else {
            document.getElementById("spanErrorResister").innerText = data
            setTimeout(()=>{
            window.location.reload()
            }, 2000)
        }

    }
    )
}