// Title: Simple replace using node

const { readFile, writeFile } = require("fs").promises;

async function replace() {
  let data = await readFile("./README.md", "utf8");
  data = data.replace("apple", "orange");
  await writeFile("./README.md", data, "utf8");
}

replace();
