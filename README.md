# copy 适应自己的项目 加入自定义schema定义

* 新增 `Intercept` 函数 可以传入拦截生成 返回false则不会生成Schema
    * 用于比如新增时或修改时 跳过某些字段的生成
* 新增 `OpenMetaData` bool 注入meta数据 目前用于标记数字的精确类型
    * 会额外生成一个meta_data{kind:准确golang kind 字符串}的map
    * 当然了 你也可以手动设置 MetaData注入你想要注入的其他内容
* 新增 `TagMapper map[string]TagMapperFunc` 用来自定义tag映射
    * 快捷方法 `AddTagSetMapper` 用来处理简单的tag赋值
        * 例如指定tag为 `comment` 赋值的字段为 `Title` 则会把Title的值赋值为 comment标签设定的内容
        * comment="someLike" 最终会设置schema的Title为 someLike
        * 是字段名 大写开头
    * 通用方法 `AddTagMapper` 自定义设置tag以及对应的处理方法函数