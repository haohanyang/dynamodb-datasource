import dayjs from "dayjs";
import timezone from "dayjs/plugin/timezone";
import advancedFormat from "dayjs/plugin/advancedFormat";
import utc from "dayjs/plugin/utc";


dayjs.extend(utc);
dayjs.extend(timezone);
dayjs.extend(advancedFormat);

export function formatRefTime(format: string) {
    // Golang's reference time: Mon, 02 Jan 2006 15:04:05 MST
    const refTime = dayjs.tz(1136239445999, "America/Edmonton");
    return refTime.format(format);
}
