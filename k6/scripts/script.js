import http from "k6/http";
import { sleep } from "k6";

const dbs = ["db1", "db2", "db3", "db4", "db5"];

function getRandomInt(max) {
  return Math.floor(Math.random() * max);
}

function makeid(length) {
  var result = "";
  var characters =
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  var charactersLength = characters.length;
  for (var i = 0; i < length; i++) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength));
  }
  return result;
}

export const options = {
  vus: 20,
  duration: "15s",
};

export default function () {
  var data = [];
  // Perform SET
  for (var i = 0; i < 1; i++) {
    const s = `db/${dbs[getRandomInt(5)]}/${makeid(8)}`;
    data.push(s);
    http.post(`http://host.docker.internal:9090/${s}`, s);
  }

  // Perform GET
  data.forEach((d) => http.get(`http://host.docker.internal:9090/${d}`));
}
