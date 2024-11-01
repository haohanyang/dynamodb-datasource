const dayjs = require("dayjs");
const timezone = require("dayjs/plugin/timezone");
const advancedFormat = require("dayjs/plugin/advancedFormat");
const utc = require("dayjs/plugin/utc");
const { exit } = require("process");

dayjs.extend(utc);
dayjs.extend(timezone);
dayjs.extend(advancedFormat);


if (process.argv.length < 3) {
    console.error("Please enter format");
    exit(1);
}

const format = process.argv[2].trim();
const refTime = dayjs.tz(1136239445999, "America/Edmonton");

console.log(refTime.format(format));
