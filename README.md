# homework
## week2
dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层？  
1）在我们公司目前的实现架构中，db操作数据时会对接多个上层业务，所以在这里是需要Wrap一个Error, 
让上层业务更加清晰地判断是什么原因导致了error  
2）如果在同一个模块内进行db的操作，对于sql.ErrNoRows可以使用wrapWithMsg添加一些query相关的信息
向上抛出，对于其他err可以直接wrap.
 
