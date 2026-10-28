import importlib.util
from pathlib import Path
import unittest


SCRIPT = Path(__file__).with_name("diagnose_operation_draw_order.py")
SPEC = importlib.util.spec_from_file_location("draw_order_diagnostic", SCRIPT)
diagnostic = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(diagnostic)


class DiagnosticFilterTest(unittest.TestCase):
    def setUp(self):
        self.old = {
            "services": [{"id": "api"}, {"id": "product-catalog"}],
            "edges": [{"id": "api-catalog", "from": "api", "to": "product-catalog"}],
        }
        self.new = {
            "services": [{"id": "api"}, {"id": "database:one"}, {"id": "product-catalog"}],
            "edges": [{"id": "api-catalog", "from": "api", "to": "product-catalog"},
                      {"id": "catalog-db", "from": "product-catalog", "to": "database:one"}],
        }

    def test_only_new_database_relation_is_removed_without_mutating_input(self):
        database_id, edge_id = diagnostic.added_database(self.old, self.new)
        filtered = diagnostic.filter_model(self.new, database_id, edge_id)
        self.assertEqual([item["id"] for item in filtered["services"]], ["api", "product-catalog"])
        self.assertEqual([item["id"] for item in filtered["edges"]], ["api-catalog"])
        self.assertEqual(len(self.new["services"]), 3)

    def test_rejects_another_incident_edge(self):
        self.new["edges"].append({"id": "api-db", "from": "api", "to": "database:one"})
        with self.assertRaises(ValueError):
            diagnostic.added_database(self.old, self.new)
        with self.assertRaises(ValueError):
            diagnostic.filter_model(self.new, "database:one", "catalog-db")

    def test_rejects_unrelated_added_service(self):
        self.new["services"][1]["id"] = "unrelated"
        with self.assertRaises(ValueError):
            diagnostic.added_database(self.old, self.new)


if __name__ == "__main__":
    unittest.main()
