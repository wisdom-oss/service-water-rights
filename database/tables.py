import geoalchemy2
import sqlalchemy

import database

geospatial_metadata = sqlalchemy.MetaData(schema="geodata")
water_right_metadata = sqlalchemy.MetaData(schema="nlwkn_water_rights")

shapes = sqlalchemy.Table(
    "shapes",
    geospatial_metadata,
    sqlalchemy.Column("id", sqlalchemy.Integer, primary_key=True, autoincrement=True),
    sqlalchemy.Column("name", sqlalchemy.Text),
    sqlalchemy.Column("key", sqlalchemy.Text),
    sqlalchemy.Column("nuts_key", sqlalchemy.Text),
    sqlalchemy.Column("geom", geoalchemy2.Geometry("Multipolygon")),
)

locations = sqlalchemy.Table(
    "usage_locations",
    water_right_metadata,
    sqlalchemy.Column("id", sqlalchemy.Integer, primary_key=True, autoincrement=True),
    sqlalchemy.Column("water_right", sqlalchemy.Integer, sqlalchemy.ForeignKey("water_rights.no")),
    sqlalchemy.Column("active", sqlalchemy.Boolean),
    sqlalchemy.Column("location", geoalchemy2.Geometry("Point", 25832)),
    sqlalchemy.Column("real", sqlalchemy.Boolean),
)


def init():
    geospatial_metadata.create_all(bind=database.engine)
    water_right_metadata.create_all(bind=database.engine)
