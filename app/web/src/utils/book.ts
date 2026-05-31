// 字数格式化，123万
export function formatWordCount(count: number) {
  if (count >= 10000) {
    return `${(count / 10000).toFixed(1)}万`
  } else {
    return count.toLocaleString()
  }
}

export function formatTime(time: string) {
  return new Date(time).toLocaleString()
}
