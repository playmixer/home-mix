const API = "http://localhost:8090/api/v0"
const tableData = document.getElementById("table_data")


async function initDatas() {
    const url = API+"/ping"
    let res = await fetch(url)
    if (res.ok) {
        let json = await res.json()
        let data = json.sort((a, b) => {
            arr = a.IP.split(".")
            brr = b.IP.split(".")
            return Number(arr[3]) < Number(brr[3]) ? -1 : 1
        })
        for (var i = 0; i < data.length; i++){
            row = tableData.insertRow(i)
            row.id = "ip"+json[i].IP
            row.insertCell(0).innerHTML = json[i].IP
            row.insertCell(1).innerHTML = json[i].Name
            row.insertCell(2).innerHTML = json[i].Mac
            row.insertCell(3).innerHTML = json[i].Online
        }

    } else {
        console.error("HTTP error: "+res.status)
    }
}

window.onload = initDatas