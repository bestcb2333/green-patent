import dayjs from "dayjs";

export function formatDate(_1, _2, time: string) {
  return dayjs(time).format('MM月DD日 HH:mm')
}
