import dayjs from "dayjs";

export function formatDate(_1:any, _2:any, time: string) {
  return dayjs(time).format('MM月DD日 HH:mm')
}
