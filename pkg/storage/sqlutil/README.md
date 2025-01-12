# sqlutil

## How to use

1. prepare a struct that represents a table in the database

    ```go
    // the primary key of the table is (id, id2)
    type User struct {
        ID       int    `column:"id" pk:"true"`
        // if no tag, the field name in snake case will be used as column name
        // in this case, the column name is `id2`
        ID2      int    `pk:"true"`
        Name     string `column:"name"`
        Age      int    `column:"age"`
        // in this case, the column name is `is_female`
        IsFemale bool   
    }
    ```

    Supported tags:
        - `column`: the column name in the database
        - `pk`: primary key, value is `true` or other string
        - `generated`: generated column, value is `true` or other string, if the column is generated, it will be ignored when inserting

    NOTE: if no `pk` tag is specified, the ID or Id field will be used as the primary key


2. create an `AppDatabaseContext`, and then use public function in this package
   to interact with the database



