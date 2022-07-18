import geoalchemy2
import sqlalchemy
import sqlalchemy.dialects.postgresql
import sqlalchemy_utils.types
import sqlalchemy.orm

import database

geospatial_metadata = sqlalchemy.MetaData(schema="geodata")
water_right_metadata = sqlalchemy.MetaData(schema="nlwkn_water_rights")
orm_base = sqlalchemy.orm.declarative_base(metadata=water_right_metadata)
sqlalchemy_utils.force_auto_coercion()
# %% Custom Composite Types
numeric_keyed_name = sqlalchemy_utils.types.CompositeType(
    "numeric_keyed_name",
    [sqlalchemy.Column("key", sqlalchemy.Integer, key="key"), sqlalchemy.Column("name", sqlalchemy.Text, key="name")],
)
interval_rate = sqlalchemy_utils.types.CompositeType(
    "interval_rate",
    [
        sqlalchemy.Column("amount", sqlalchemy.Integer, key="amount"),
        sqlalchemy.Column("unit", sqlalchemy.Text, key="unit"),
        sqlalchemy.Column("duration", sqlalchemy.Interval, key="duration"),
    ],
)
rate = sqlalchemy_utils.types.CompositeType(
    "rate",
    [
        sqlalchemy.Column("amount", sqlalchemy.Integer, key="amount"),
        sqlalchemy.Column("unit", sqlalchemy.Text, key="unit"),
    ],
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


class E_UsageLocation(orm_base):
    __tablename__ = "e_usage_locations"

    id = sqlalchemy.Column("id", sqlalchemy.Integer, primary_key=True, autoincrement=True)
    water_right = sqlalchemy.Column("water_right", sqlalchemy.Integer, sqlalchemy.ForeignKey("water_rights.no"))
    name = sqlalchemy.Column("name", sqlalchemy.Text)
    no = sqlalchemy.Column("no", sqlalchemy.Integer)
    active = sqlalchemy.Column("active", sqlalchemy.Boolean)
    location = sqlalchemy.Column("location", geoalchemy2.Geometry("Point", 25832))
    basin_no = sqlalchemy.Column("basin_no", numeric_keyed_name)
    county = sqlalchemy.Column("county", sqlalchemy.Text)
    eu_survey_area = sqlalchemy.Column("eu_survey_area", numeric_keyed_name)
    field = sqlalchemy.Column("field", sqlalchemy.Integer)
    groundwater_volume = sqlalchemy.Column("groundwater_volume", sqlalchemy.Text)
    legal_scope = sqlalchemy.Column("legal_scope", sqlalchemy.Text)
    local_sub_district = sqlalchemy.Column("local_sub_district", sqlalchemy.Text)
    maintenance_association = sqlalchemy.Column("maintenance_association", numeric_keyed_name)
    municipal_area = sqlalchemy.Column("municipal_area", numeric_keyed_name)
    plot = sqlalchemy.Column("plot", sqlalchemy.Text)
    real = sqlalchemy.Column("real", sqlalchemy.Boolean)
    rivershed = sqlalchemy.Column("rivershed", sqlalchemy.Text)
    serial_no = sqlalchemy.Column("serial_no", sqlalchemy.Text)
    top_map_1_25000 = sqlalchemy.Column("top_map_1_25000", numeric_keyed_name)
    water_body = sqlalchemy.Column("water_body", sqlalchemy.Text)
    flood_area = sqlalchemy.Column("flood_area", sqlalchemy.Text)
    water_protection_area = sqlalchemy.Column("water_protection_area", sqlalchemy.Text)
    withdrawal_rate = sqlalchemy.Column("withdrawal_rate", sqlalchemy.dialects.postgresql.ARRAY(interval_rate))
    fluid_discharge = sqlalchemy.Column("fluid_discharge", sqlalchemy.dialects.postgresql.ARRAY(interval_rate))
    irrigation_area = sqlalchemy.Column("irrigation_area", rate)
    rain_supplement = sqlalchemy.Column("rain_supplement", sqlalchemy.dialects.postgresql.ARRAY(interval_rate))


class WaterRight(orm_base):
    __tablename__ = "water_rights"
    id = sqlalchemy.Column("id", sqlalchemy.Integer, primary_key=True, autoincrement=True)
    no = sqlalchemy.Column("no", sqlalchemy.Integer)
    ext_id = sqlalchemy.Column("ext_id", sqlalchemy.Text)
    file_ref = sqlalchemy.Column("file_ref", sqlalchemy.Text)
    legal_title = sqlalchemy.Column("legal_title", sqlalchemy.Text)
    state = sqlalchemy.Column("state", sqlalchemy.Enum("aktiv", "inaktiv", "Wasserbuchblatt", name="water_right_state"))
    subject = sqlalchemy.Column("subject", sqlalchemy.Text)
    address = sqlalchemy.Column("address", sqlalchemy.Text)
    annotation = sqlalchemy.Column("annotation", sqlalchemy.Text)
    bailee = sqlalchemy.Column("bailee", sqlalchemy.Text)
    date_of_change = sqlalchemy.Column("date_of_change", sqlalchemy.dialects.postgresql.TIMESTAMP)
    valid = sqlalchemy.Column("valid", sqlalchemy.dialects.postgresql.DATERANGE)
    granting_authority = sqlalchemy.Column("granting_authority", sqlalchemy.Text)
    registering_authority = sqlalchemy.Column("registering_authority", sqlalchemy.Text)
    water_authority = sqlalchemy.Column("water_authority", sqlalchemy.Text)


def init():
    sqlalchemy_utils.register_composites(database.engine.connect())
    geospatial_metadata.create_all(bind=database.engine)
    water_right_metadata.create_all(bind=database.engine)
