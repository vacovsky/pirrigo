package main

//    sqlStr = """
//        SELECT DISTINCT dripnodes.sid,
//            SUM((duration / 60 )) as runmins,
//            (SELECT sum((gph * count)) as totalgph from dripnodes where dripnodes.sid=history.sid) as totalgph,
//            stations.notes
//        FROM history
//        INNER JOIN dripnodes ON dripnodes.sid=history.sid
//        INNER JOIN stations ON stations.id=history.sid
//            WHERE starttime >= (CURRENT_DATE - INTERVAL 30 DAY)
//            GROUP BY dripnodes.sid
//            ORDER BY dripnodes.sid ASC;
//            """
//    results = {'water_usage': []}
//    for d in sqlConn.read(sqlStr):
//        results['water_usage'].append(
//            {
//                'sid': int(d[0]),
//                'notes': str(d[3]),
//                'run_mins': int(d[1]),
//                'total_gph': float(d[2]),
//                'usage_last_30': float((d[1] / 60) * d[2])
//            }
//        )
