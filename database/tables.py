import geoalchemy2
import sqlalchemy
import sqlalchemy.dialects.postgresql
import sqlalchemy_utils.types

import database

geospatial_metadata = sqlalchemy.MetaData(schema="geodata")
water_right_metadata = sqlalchemy.MetaData(schema="nlwkn_water_rights")
sqlalchemy_utils.force_auto_coercion()
# %% Custom Composite Types
numeric_keyed_name = sqlalchemy_utils.types.CompositeType(
    "numeric_keyed_name", [sqlalchemy.Column("key", sqlalchemy.Integer), sqlalchemy.Column("name", sqlalchemy.Text)]
)
interval_rate = sqlalchemy_utils.types.CompositeType(
    "interval_rate",
    [
        sqlalchemy.Column("amount", sqlalchemy.Integer),
        sqlalchemy.Column("unit", sqlalchemy.Text),
        sqlalchemy.Column("duration", sqlalchemy.Interval),
    ],
)
rate = sqlalchemy_utils.types.CompositeType(
    "rate",
    [sqlalchemy.Column("amount", sqlalchemy.Integer), sqlalchemy.Column("unit", sqlalchemy.Text)],
)
# %% Tables

shapes = sqlalchemy.Table(
    "shapes",
    geospatial_metadata,
    sqlalchemy.Column("id", sqlalchemy.Integer, primary_key=True, autoincrement=True),
    sqlalchemy.Column("name", sqlalchemy.Text),
    sqlalchemy.Column("key", sqlalchemy.Text),
    sqlalchemy.Column("nuts_key", sqlalchemy.Text),
    sqlalchemy.Column("geom", geoalchemy2.Geometry("Multipolygon")),
)

e_usage_locations = sqlalchemy.Table(
    "e_usage_locations",
    water_right_metadata,
    sqlalchemy.Column("id", sqlalchemy.Integer, primary_key=True, autoincrement=True),
    sqlalchemy.Column("water_right", sqlalchemy.Integer, sqlalchemy.ForeignKey("water_rights.no")),
    sqlalchemy.Column("name", sqlalchemy.Text),
    sqlalchemy.Column("no", sqlalchemy.Integer),
    sqlalchemy.Column("active", sqlalchemy.Boolean),
    sqlalchemy.Column("location", geoalchemy2.Geometry("Point", 25832)),
    sqlalchemy.Column("basin_no", numeric_keyed_name),
    sqlalchemy.Column("county", sqlalchemy.Text),
    sqlalchemy.Column("eu_survey_area", numeric_keyed_name),
    sqlalchemy.Column("field", sqlalchemy.Integer),
    sqlalchemy.Column("groundwater_volume", sqlalchemy.Text),
    sqlalchemy.Column("legal_scope", sqlalchemy.Text),
    sqlalchemy.Column("local_sub_district", sqlalchemy.Text),
    sqlalchemy.Column("maintenance_association", numeric_keyed_name),
    sqlalchemy.Column("municipal_area", numeric_keyed_name),
    sqlalchemy.Column("plot", sqlalchemy.Text),
    sqlalchemy.Column("real", sqlalchemy.Boolean),
    sqlalchemy.Column("rivershed", sqlalchemy.Text),
    sqlalchemy.Column("serial_no", sqlalchemy.Text),
    sqlalchemy.Column("top_map_1_25000", numeric_keyed_name),
    sqlalchemy.Column("water_body", sqlalchemy.Text),
    sqlalchemy.Column("flood_area", sqlalchemy.Text),
    sqlalchemy.Column("water_protection_area", sqlalchemy.Text),
    sqlalchemy.Column("withdrawal_rate", sqlalchemy.dialects.postgresql.ARRAY(interval_rate)),
    sqlalchemy.Column("fluid_discharge", sqlalchemy.dialects.postgresql.ARRAY(interval_rate)),
    sqlalchemy.Column("irrigation_area", rate),
    sqlalchemy.Column("rain_supplement", sqlalchemy.dialects.postgresql.ARRAY(interval_rate)),
)

water_rights = sqlalchemy.Table(
    "water_rights",
    water_right_metadata,
    sqlalchemy.Column("id", sqlalchemy.Integer, primary_key=True, autoincrement=True),
    sqlalchemy.Column("no", sqlalchemy.Integer),
    sqlalchemy.Column("ext_id", sqlalchemy.Text),
    sqlalchemy.Column("file_ref", sqlalchemy.Text),
    sqlalchemy.Column("legal_title", sqlalchemy.Text),
    sqlalchemy.Column("state", sqlalchemy.dialects.postgresql.ENUM("aktiv", "inaktiv", "Wasserbuchblatt")),
    sqlalchemy.Column("subject", sqlalchemy.Text),
    sqlalchemy.Column("address", sqlalchemy.Text),
    sqlalchemy.Column("annotation", sqlalchemy.Text),
    sqlalchemy.Column("bailee", sqlalchemy.Text),
    sqlalchemy.Column("date_of_change", sqlalchemy.dialects.postgresql.TIMESTAMP),
    sqlalchemy.Column("valid", sqlalchemy.dialects.postgresql.DATERANGE),
    sqlalchemy.Column("granting_authority", sqlalchemy.Text),
    sqlalchemy.Column("registering_authority", sqlalchemy.Text),
    sqlalchemy.Column("water_authority", sqlalchemy.Text),
)


def init():
    geospatial_metadata.create_all(bind=database.engine)
    water_right_metadata.create_all(bind=database.engine)
