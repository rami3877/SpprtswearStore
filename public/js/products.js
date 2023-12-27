function Getproducts(kind , container , id , callfunction){
  
    fetch(`/product?kind=${kind}&container=${container}&id=${id}`, {
        method: "GET",
    }).then(Response => Response.json()).then(data => {
      callfunction(data)
    })
}