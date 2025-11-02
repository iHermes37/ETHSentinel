

type MySqlPool=Pool<ConnectionManager<MysqlConnection>>;

pub struct DbSessionManager {
    pool: MySqlPool,
}

impl DbSessionManager{
    pub fn new(database_url: &str)->Self{

        //Diesel/r2d2 提供的 MySQL 连接管理器，用来管理数据库连接
        let manager=ConnectionManager::<MysqlConnection>::new(database_url);

        //创建一个 连接池构建器,用上面创建的 manager 构建连接池,
        //如果构建失败，直接 panic 并打印 "Failed to create pool"
        let pool=r2d2::Pool::builder().build(manager)
            .expect("Failed to create pool");

        DbSessionManager(pool)
    }

    pub fn get_connection(&self) -> Result<PooledConnection<ConnectionManager<MysqlConnection>>, Error> {
        self.pool.get()
    }
}