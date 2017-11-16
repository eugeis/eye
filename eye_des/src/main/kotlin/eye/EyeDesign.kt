package eye

import ee.design.*
import ee.lang.*

object Eye : Comp() {
    object Service : CompilationUnit({ ifc(true) }) {
        val name = propS()

    }

    object Query : CompilationUnit({ ifc(true) }) {
        val info = propS()

    }

    object Check : CompilationUnit({ ifc(true) }) {
        val info = propS()

    }

    object QueryResult : Values({ ifc(true) }) {
        val info = propS()

    }

    object MySqlService : Entity({ superUnit(Service) }) {
    }

    object WebService : Entity({ superUnit(Service) }) {
    }

    object ValidationRequest : Basic() {
        val query = propS()
        val regExpr = propS()
        val evalExpr = propS()
        val all = propB()
    }

    object ExportRequest : Basic() {
        val query = propS()
        val evalExpr = propS()
    }

    object CommandRequest : Basic() {
    }
}