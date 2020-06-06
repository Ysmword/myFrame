dgw postgres://postgres:postgres@119.23.67.53:5432/gzm --schema=public --package=models --output=./models/models.go --template=./tmpl/struct.tmpl --typemap=./tmpl/typemap.toml 



    # dgw --schema=public \ 
    #  schema 可以在创建表的时候可以看到，默认是使用public
    # --package=models \
    # package 放置的包
    # --output=./models/models.go \
    # output 生成指定文件
    # --template=./tmpl/struct.tmpl \
    # template 指定模板
    # --typemap=./tmpl/typemap.toml \
    # typemap列类型和go类型映射文件路径
    # postgres://postgres:postgres@119.23.67.53:5432/gzm
    # postgres为用户名，postgres为密码119.23.67.53:5432顾名思义为数据库所在的地址，gzm为数据库名称