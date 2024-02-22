const API = "/api/v0"
const tableData = document.getElementById("table_data")


async function initDatas() {
    const url = API+"/ping"
    let res = await fetch(url)
    if (res.ok) {
        let json = await res.json()
        let data = json.sort((a, b) => {
            arr = a.ip.split(".")
            brr = b.ip.split(".")
            return Number(arr[3]) < Number(brr[3]) ? -1 : 1
        })
        for (var i = 0; i < data.length; i++){
            row = tableData.insertRow(i)
            row.id = "ip"+json[i].ip
            row.insertCell(0).innerHTML = i + 1
            row.insertCell(1).innerHTML = json[i].ip
            row.insertCell(2).innerHTML = json[i].name
            row.insertCell(3).innerHTML = json[i].mac
            rowOnline = row.insertCell(4)

            
            if (String(json[i].online).trim() == "true"){
                rowOnline.classList.add("online")
                rowOnline.innerHTML = "Ok"
            } else {
                rowOnline.classList.add("offline")
                rowOnline.innerHTML = "Disabled"
            }
        }

    } else {
        console.error("HTTP error: "+res.status)
    }
}

window.onload = initDatas