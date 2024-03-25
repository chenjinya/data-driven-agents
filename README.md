# data-driven-agents

一种有别于 LangChain 的 Agent 流程调度方式

## Roadmap

```golang
roadmap := pipline.RoadPath{
    ID: "entry",
    Next: []pipline.RoadPath{
        {
            ID: "number is odd ?",
            Next: []pipline.RoadPath{
                {
                    ID: "judgement",
                },
            },
        },
    },
}
```