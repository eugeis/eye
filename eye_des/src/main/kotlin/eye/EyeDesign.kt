package eye

import ee.design.*
import ee.lang.*

object Eye : Comp({ namespace("eye") }) {
    object Shared : Module() {
        object Service : CompilationUnit({ ifc(true) }) {
            val name = propS()
            val init = op { errorHandling(false) }
            val close = op()
            val ping = op { errorHandling(false) }

            val newCheck = op(p("req", ValidationRequest)) { ret(Check).errorHandling(false) }
            val newExporter = op(p("req", ExportRequest)) { ret(Exporter).errorHandling(false) }
            val newExecutor = op(p("req", CommandRequest)) { ret(Executor).errorHandling(false) }
        }

        object Query : CompilationUnit({ ifc(true) }) {
            val info = propS()
            val query = op { ret(n.List.G(QueryResult)).errorHandling(false) }
        }

        object Check : CompilationUnit({ ifc(true).superUnit(Query) }) {
            val validate = op { errorHandling(false) }
        }

        object Exporter : CompilationUnit({ ifc(true) }) {
            val info = propS()
            val export = op(p("params", n.Map)) { errorHandling(false) }
        }

        object Executor : CompilationUnit({ ifc(true) }) {
            val info = propS()
            val execute = op(p("params", n.Map)) { errorHandling(false) }
        }

        object QueryResult : Values({ ifc(true) }) {
            val info = propS()

        }

        object ExportRequest : Basic() {
            val query = propS()
            val evalExpr = propS()
        }

        object ValidationRequest : Basic({ superUnit(ExportRequest) }) {
            val regExpr = propS()
            val all = propB()
        }

        object CommandRequest : Basic() {}

        object ServiceFactory : CompilationUnit({ ifc(true) }) {
            val find = op(p("name", n.String)) { ret(Service).errorHandling(false) }
            val close = op()
        }

        object MultiCheck : Controller({ superUnit(Shared.Check) }) {
            val info = propS()
            val queries = prop(n.List.G(Shared.Query))
            val all = propB()
            val onlyRunning = propB()

            val execute = op(p("params", n.Map)) { errorHandling(false) }
        }
    }

    object MySql : Module() {
        object MySqlService : Entity({ superUnit(Shared.Service) }) {}
    }

    object Http : Module() {
        object HttpService : Entity({ superUnit(Shared.Service) }) {}
    }

    object Process : Module() {
        object ProcessService : Entity({ superUnit(Shared.Service) }) {}
    }

    object FileSystem : Module() {
        object FileSystemService : Entity({ superUnit(Shared.Service) }) {}
    }

    object Elastic : Module() {
        object ElasticService : Entity({ superUnit(Shared.Service) }) {}
    }
}