fetch("/api/login", {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
  },
  body: `{"username" :"d"}`,
}).then(Response => Response.json()).then( data =>{
	 console.log(data)
}
);
