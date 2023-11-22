package main  
  
import (  
    "fmt"  
    "math/rand"  
    "sync"  
    "time"  
)  
  
const (  
    n     = 4 // 系统中副本的数量  
    k     = 3 // 系统中正确副本的数量  
    id    = 0 // 本机的ID  
    delay = 2 * time.Second // 延迟时间  
)  
  
var wg sync.WaitGroup  
  
func main() {  
    // 初始化本地状态  
    local := make([]int, n)  
    for i := 0; i < n; i++ {  
        local[i] = -1  
    }  
    local[id] = 0  
  
    // 定义回调函数，用于处理接收到的提议  
    delivered := func(value int) {  
        if local[value] == -1 {  
            local[value] = 1  
        } else {  
            local[value] = 0  
        }  
    }  
  
    // 定义广播函数，用于将提议广播给其他副本  
    broadcast := func(value int) {  
        for i := 0; i < n; i++ {  
            if i != id {  
                go func(id int) {  
                    conn, err := fmt.Dial("tcp", fmt.Sprintf("replica%d:1234", id))  
                    if err != nil {  
                        fmt.Println(err)  
                        return  
                    }  
                    defer conn.Close()  
                    fmt.Fprintf(conn, "delivered:%d\n", value)  
                }(i)  
            }  
        }  
    }  
  
    // 进入WRBC阶段，等待其他副本提议值并广播给其他副本自己收到的值  
    for i := 0; i < n-k+1; i++ {  
        wg.Add(1)  
        go func() {  
            defer wg.Done()  
            for j := 0; j < k; j++ {  
                time.Sleep(delay) // 等待一段时间，模拟延迟传输  
                value := rand.Intn(n) // 生成一个随机值作为提议值  
                broadcast(value) // 将提议值广播给其他副本  
                delivered(value) // 处理收到的提议值，更新本地状态  
            }  
        }()  
    }  
    wg.Wait() // 等待所有goroutine执行完毕，确保所有副本都收到的提议值都处理完毕  

  // 进入MVBA阶段  
  for {  
    // 选择一个输出值  
      var output int  
      for {  
          output = -1  
          for i := 0; i < n; i++ {  
              if local[i] == 1 {  
                  output = i  
                  break  
              }  
          }  
          if output == -1 {  
              // 没有正确的副本选择输出值，重新进入MVBA阶段  
              continue  
          }  
          // 检查其他副本的选择值是否与自己相同  
          for i := 0; i < n; i++ {  
              if i != id && local[i] == 0 && output != i {  
                  // 其他副本的选择值与自己不同，重新进入MVBA阶段  
                  continue outer  
              }  
          }  
          break  
      }  
      // 广播自己的选择值  
      broadcast(output)  
      // 处理收到的其他副本的选择值，更新本地状态  
      delivered(output)  
      // 检查自己是否是正确副本，如果是则输出最终结果，否则输出收到的其他副本的选择值  
      if local[output] == 1 {  
          fmt.Printf("Replica %d: Output: %d\n", id, output)  
      } else {  
          fmt.Printf("Replica %d: Output: %d (Received from Replica %d)\n", id, output, output)  
      }  
  }
}
